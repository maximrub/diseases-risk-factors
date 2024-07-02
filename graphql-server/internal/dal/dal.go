package dal

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/bsm/redislock"
	"github.com/google/uuid"
	"github.com/hashicorp/go-set"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/auth"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal/kegg"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal/pubmed"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/graph/model"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/utils"
)

type DAL struct {
	dbClient                      *mongo.Client
	redisClient                   *redis.Client
	locker                        *redislock.Client
	storageAccountName            string
	storageAccountUrl             string
	storageAccountKey             string
	diseasesCollection            *mongo.Collection
	articlesCollection            *mongo.Collection
	qasCollection                 *mongo.Collection
	classificationItemsCollection *mongo.Collection
	riskFactorsCollection         *mongo.Collection
	articleSearcher               ArticleSearcher
	diseasesFetcher               DiseasesFetcher
}

type ArticleSearcher interface {
	Search(query string, limit int) ([]int, error)
	Fetch(id int) (string, error)
}

type DiseasesFetcher interface {
	List() ([]string, error)
	Fetch(id string) (*model.Disease, error)
}

func NewDal(ctx context.Context) (*DAL, error) {
	connectionUri := os.Getenv("DB_CONNECTION_URI")
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, utils.WrapError(err, "error connecting to DB")
	}

	// Check the connection
	err = dbClient.Ping(ctx, nil)
	if err != nil {
		return nil, utils.WrapError(err, "error ping the DB")
	}

	database := dbClient.Database(os.Getenv("DATABASE_NAME"))

	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		Password:     redisPassword,
		TLSConfig:    &tls.Config{MinVersion: tls.VersionTLS12},
		WriteTimeout: 5 * time.Second,
	})

	storageAccountName := os.Getenv("STORAGE_ACCOUNT_NAME")
	storageAccountUrl := fmt.Sprintf("https://%s.blob.core.windows.net", storageAccountName)

	return &DAL{
		dbClient:                      dbClient,
		redisClient:                   redisClient,
		locker:                        redislock.New(redisClient),
		storageAccountName:            storageAccountName,
		storageAccountUrl:             storageAccountUrl,
		storageAccountKey:             os.Getenv("STORAGE_ACCOUNT_KEY"),
		diseasesCollection:            database.Collection("diseases"),
		articlesCollection:            database.Collection("articles"),
		qasCollection:                 database.Collection("qas"),
		classificationItemsCollection: database.Collection("classificationItems"),
		riskFactorsCollection:         database.Collection("riskFactors"),
		articleSearcher:               pubmed.NewEntrezSearcher(),
		diseasesFetcher:               kegg.NewDiseasesFetcher(),
	}, nil
}

func (d *DAL) GetDiseases(ctx context.Context) ([]*model.Disease, error) {
	var diseases []*Disease
	cursor, err := d.diseasesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding diseases in DB")
	}

	if err = cursor.All(ctx, &diseases); err != nil {
		return nil, utils.WrapErrorf(err, "error loading diseases from DB")
	}

	var gqlDiseases []*model.Disease
	for _, disease := range diseases {
		gqlDisease := convertDiseaseToGql(disease)
		gqlDiseases = append(gqlDiseases, gqlDisease)
	}

	return gqlDiseases, nil
}

func (d *DAL) CreateQA(ctx context.Context, input model.CreateQAInput) (*model.QAPayload, error) {
	if !auth.HasPermission(ctx, "write:qas") {
		return nil, auth.NotAuthorizedError
	}

	disease, err := d.getDisease(input.DiseaseID)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting disease [%s]", input.DiseaseID)
	}

	article, err := d.getArticle(input.ArticleID)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting article [%s]", input.ArticleID)
	}

	qa := &QA{
		ID:      uuid.NewString(),
		Disease: disease,
		Article: article,
	}

	for _, question := range input.Questions {
		qaQuestion := &Question{
			ID:   uuid.NewString(),
			Text: question.Text,
		}

		for _, answer := range question.Answers {
			qaQuestion.Answers = append(qaQuestion.Answers, &Answer{
				AnswerStart: answer.AnswerStart,
				Text:        answer.Text,
			})
		}

		qa.Questions = append(qa.Questions, qaQuestion)
	}

	_, err = d.qasCollection.InsertOne(ctx, qa)
	if err != nil {
		return nil, utils.WrapError(err, "error inserting qa to DB")
	}

	return &model.QAPayload{Qa: convertQAToGql(qa)}, nil
}

func (d *DAL) DeleteQA(ctx context.Context, input model.DeleteQAInput) (*model.DeleteQAPayload, error) {
	if !auth.HasPermission(ctx, "write:qas") {
		return nil, auth.NotAuthorizedError
	}

	_, err := d.qasCollection.DeleteOne(ctx, bson.M{"_id": input.ID})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error deleting qa [%s] in DB", input.ID)
	}

	return &model.DeleteQAPayload{}, nil
}

func (d *DAL) FetchDiseases(ctx context.Context) error {
	if !auth.HasPermission(ctx, "write:diseases") {
		return auth.NotAuthorizedError
	}

	const lockKey = "FetchDiseases"

	lockCtx := log.WithContext(context.Background()).WithField("lock", lockKey).Context
	logger := log.WithContext(lockCtx).Logger
	lock, err := d.locker.Obtain(lockCtx, lockKey, 30*time.Second, nil)
	if err != nil {
		return utils.WrapErrorf(err, "error obtain [%s] lock", lockKey)
	}

	go func() {
		defer lock.Release(lockCtx)

		diseasesIDs, err2 := d.diseasesFetcher.List()
		if err2 != nil {
			logger.WithError(err2).Error("error listing diseases IDs")
			return
		}

		for _, id := range diseasesIDs {
			// Extend lock
			if err3 := lock.Refresh(lockCtx, 30*time.Second, nil); err != nil {
				logger.WithError(err3).Errorf("error extending [%s] lock", lockKey)
				return
			}

			disease, err3 := d.diseasesFetcher.Fetch(id)
			if err3 != nil {
				logger.WithError(err3).Errorf("error getting disease [%s]", id)
				continue
			}

			diseaseModel := convertDiseaseToDB(disease)
			_, err3 = d.diseasesCollection.ReplaceOne(lockCtx, bson.M{"_id": diseaseModel.ID}, diseaseModel, options.Replace().SetUpsert(true))
			if err3 != nil {
				logger.WithError(err3).Errorf("error inserting disease [%s] to DB", id)
				return
			}
		}
	}()

	return nil
}

func (d *DAL) CreateDisease(ctx context.Context, input model.CreateDiseaseInput) (*model.DiseasePayload, error) {
	if !auth.HasPermission(ctx, "write:diseases") {
		return nil, auth.NotAuthorizedError
	}

	disease := &Disease{
		ID:          fmt.Sprintf("M%s", time.Now().UTC().Format("2006_01_02_15_04_05")),
		Names:       input.Names,
		Category:    input.Category,
		Description: input.Description,
	}

	if input.DbLinks != nil {
		disease.DbLinks = &DiseaseDBLinks{
			Icd10: input.DbLinks.Icd10,
			Icd11: input.DbLinks.Icd11,
			Mesh:  input.DbLinks.Mesh,
		}
	}

	_, err := d.diseasesCollection.InsertOne(ctx, disease)
	if err != nil {
		return nil, utils.WrapError(err, "error inserting disease to DB")
	}

	return &model.DiseasePayload{Disease: convertDiseaseToGql(disease)}, nil
}

func (d *DAL) UpdateQA(ctx context.Context, input model.UpdateQAInput) (*model.UpdateQAPayload, error) {
	if !auth.HasPermission(ctx, "write:qas") {
		return nil, auth.NotAuthorizedError
	}

	qa, err := d.getQA(input.QaID)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting qa [%s]", input.QaID)
	}

	if input.Patch == nil {
		return &model.UpdateQAPayload{Qa: convertQAToGql(qa)}, nil
	}

	qa.Questions = make([]*Question, 0, len(input.Patch.Questions))
	for _, question := range input.Patch.Questions {
		qaQuestion := &Question{
			ID:   uuid.NewString(),
			Text: question.Text,
		}

		qaQuestion.Answers = make([]*Answer, 0, len(question.Answers))
		for _, answer := range question.Answers {
			qaQuestion.Answers = append(qaQuestion.Answers, &Answer{
				AnswerStart: answer.AnswerStart,
				Text:        answer.Text,
			})
		}

		qa.Questions = append(qa.Questions, qaQuestion)
	}

	_, err = d.qasCollection.ReplaceOne(context.TODO(), bson.M{"_id": qa.ID}, qa)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error updating qa [%s] in DB", qa.ID)
	}

	return &model.UpdateQAPayload{Qa: convertQAToGql(qa)}, nil
}

func (d *DAL) getDisease(id string) (*Disease, error) {
	var disease Disease
	err := d.diseasesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&disease)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting disease [%s] from DB", id)
	}

	return &disease, nil
}

func (d *DAL) getQA(id string) (*QA, error) {
	var qa QA
	err := d.qasCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&qa)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting qa [%s] from DB", id)
	}

	return &qa, nil
}

func (d *DAL) GetArticle(id string) (*model.Article, error) {
	article, err := d.getArticle(id)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting article [%s]", id)
	}

	return convertArticleToGql(article), nil
}

func (d *DAL) getArticle(id string) (*Article, error) {
	var article Article
	err := d.articlesCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&article)
	if err == nil {
		return &article, nil
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		// Error other than "no documents found", return it
		return nil, utils.WrapErrorf(err, "error getting article [%s] from DB", id)
	}

	articleID, err := strconv.Atoi(id)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error converting article ID [%s]", id)
	}

	articleText, err := d.articleSearcher.Fetch(articleID)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error fetching article [%s]", id)
	}

	if articleText == "" {
		return nil, utils.WrapErrorf(err, "error article [%s] text is empty", id)
	}

	article.ID = id
	article.Text = articleText
	_, err = d.articlesCollection.InsertOne(context.TODO(), article)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error inserting article [%s] to DB", id)
	}

	return &article, nil
}

func (d *DAL) GetQAs(ctx context.Context, diseaseID string) ([]*model.Qa, error) {
	var qas []*QA
	cursor, err := d.qasCollection.Find(ctx, bson.M{"disease._id": diseaseID})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding QAs in DB")
	}

	if err = cursor.All(ctx, &qas); err != nil {
		return nil, utils.WrapErrorf(err, "error loading QAs from DB")
	}

	var gqlQAs []*model.Qa
	for _, qa := range qas {
		gqlQAs = append(gqlQAs, convertQAToGql(qa))
	}

	return gqlQAs, nil
}

func (d *DAL) Close(ctx context.Context) error {
	err := d.dbClient.Disconnect(ctx)
	if err != nil {
		return utils.WrapErrorf(err, "error disconnecting from DB")
	}

	err = d.redisClient.Close()
	if err != nil {
		return utils.WrapErrorf(err, "error disconnecting from redis")
	}

	return nil
}

func (d *DAL) GetClassificationItems(ctx context.Context) ([]*model.ClassificationItem, error) {
	var classificationItems []*ClassificationItem
	classificationItemsCursor, err := d.classificationItemsCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding classification items in DB")
	}

	if err = classificationItemsCursor.All(ctx, &classificationItems); err != nil {
		return nil, utils.WrapErrorf(err, "error loading classification items from DB")
	}

	var qas []*QA
	qaCursor, err := d.qasCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding QAs in DB")
	}

	if err = qaCursor.All(ctx, &qas); err != nil {
		return nil, utils.WrapErrorf(err, "error loading QAs from DB")
	}

	articlesIDs := set.New[string](0)
	var gqlClassificationItems []*model.ClassificationItem

	for _, qa := range qas {
		if !articlesIDs.Contains(qa.Article.ID) {
			classificationItem := &ClassificationItem{
				ID:      qa.Article.ID,
				Label:   ClassificationLabelPositive,
				Article: qa.Article,
			}

			gqlClassificationItem := convertClassificationItemToGql(classificationItem)
			gqlClassificationItems = append(gqlClassificationItems, gqlClassificationItem)
			articlesIDs.Insert(qa.Article.ID)
		}
	}

	for _, classificationItem := range classificationItems {
		if !articlesIDs.Contains(classificationItem.Article.ID) {
			gqlClassificationItem := convertClassificationItemToGql(classificationItem)
			gqlClassificationItems = append(gqlClassificationItems, gqlClassificationItem)
			articlesIDs.Insert(classificationItem.Article.ID)
		}
	}

	return gqlClassificationItems, nil
}

func (d *DAL) GetStatistics(ctx context.Context) (*model.Statistics, error) {
	diseasesCount, err := d.diseasesCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting diseases count in DB")
	}

	var classificationItems []*ClassificationItem
	classificationItemsCursor, err := d.classificationItemsCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding classification items in DB")
	}

	if err = classificationItemsCursor.All(ctx, &classificationItems); err != nil {
		return nil, utils.WrapErrorf(err, "error loading classification items from DB")
	}

	positiveClassificationItemsCount := 0
	negativeClassificationItemsCount := 0

	for _, classificationItem := range classificationItems {
		switch classificationItem.Label {
		case ClassificationLabelPositive:
			positiveClassificationItemsCount++
		case ClassificationLabelNegative:
			negativeClassificationItemsCount++
		}
	}

	qasCount, err := d.qasCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting QAs count in DB")
	}

	cursor, err := d.qasCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error finding QAs in DB")
	}
	defer cursor.Close(ctx)

	shortestAnswer := math.MaxInt32
	longestAnswer := 0

	for cursor.Next(ctx) {
		var qa QA
		if err2 := cursor.Decode(&qa); err2 != nil {
			return nil, utils.WrapErrorf(err2, "error decoding QA from DB")
		}

		for _, question := range qa.Questions {
			for _, answer := range question.Answers {
				length := len(answer.Text)

				if length == 217 {
					log.WithField("qa", qa.ID).WithField("answer", answer.Text).Info("answer")
				}

				if length < shortestAnswer {
					shortestAnswer = length
				}

				if length > longestAnswer {
					longestAnswer = length
				}
			}
		}
	}

	articlesCount, err := d.articlesCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting articles count in DB")
	}

	return &model.Statistics{
		DiseaseCount:                          int(diseasesCount),
		ArticlesCount:                         int(articlesCount),
		QasCount:                              int(qasCount),
		PositiveClassificationItemsCount:      positiveClassificationItemsCount,
		TotalPositiveClassificationItemsCount: positiveClassificationItemsCount + int(qasCount),
		NegativeClassificationItemsCount:      negativeClassificationItemsCount,
		QasShortestAnswerLength:               shortestAnswer,
		QasLongestAnswerLength:                longestAnswer,
	}, nil
}

func (d *DAL) CreateClassificationItems(ctx context.Context, input []*model.CreateClassificationItemInput) (*model.ClassificationItemPayload, error) {
	if !auth.HasPermission(ctx, "write:classification") {
		return nil, auth.NotAuthorizedError
	}

	resp := &model.ClassificationItemPayload{
		ClassificationItems: make([]*model.ClassificationItem, 0, len(input)),
	}

	for _, item := range input {
		article, err := d.getArticle(item.ArticleID)
		if err != nil {
			log.WithError(err).WithField("articleID", item.ArticleID).Error("Error getting article")
			continue
		}

		classificationItem := &ClassificationItem{
			ID:      article.ID,
			Label:   ClassificationLabel(item.Label),
			Article: article,
		}

		_, err = d.classificationItemsCollection.InsertOne(ctx, classificationItem)
		if err != nil {
			log.WithError(err).WithField("classificationItemID", classificationItem.ID).Error("Error inserting classification item to DB")
			continue
		}

		resp.ClassificationItems = append(resp.ClassificationItems, convertClassificationItemToGql(classificationItem))
	}

	return resp, nil
}

func (d *DAL) SearchArticles(ctx context.Context, term string, limit int) ([]string, error) {
	ids, err := d.articleSearcher.Search(term, limit)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error searching for \"%s\" with [%s] limit count", term, limit)
	}

	articlesIDs := utils.Map(ids, func(v int) string {
		return strconv.Itoa(v)
	})

	return articlesIDs, nil
}

func (d *DAL) GetRiskFactors(ctx context.Context, diseaseID string) (*model.RiskFactors, error) {
	var riskFactors RiskFactors
	err := d.riskFactorsCollection.FindOne(context.TODO(), bson.M{"_id": diseaseID}).Decode(&riskFactors)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, utils.WrapErrorf(err, "error getting risk factors for disease [%s] from DB", diseaseID)
	}

	return convertRiskFactorsToGql(&riskFactors), nil
}

func (d *DAL) UpdateRiskFactors(ctx context.Context, input model.UpdateRiskFactorsInput) (*model.UpdateRiskFactorsPayload, error) {
	if !auth.HasPermission(ctx, "write:riskFactors") {
		return nil, auth.NotAuthorizedError
	}

	disease, err := d.getDisease(input.DiseaseID)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error getting disease [%s]", input.DiseaseID)
	}

	riskFactorsModel := &RiskFactors{
		ID:      disease.ID,
		Disease: disease,
	}

	for _, riskFactor := range input.RiskFactors {
		riskFactorsModel.RiskFactors = append(riskFactorsModel.RiskFactors, &RiskFactor{
			Text:        riskFactor.Text,
			Score:       riskFactor.Score,
			ArticlesIDs: riskFactor.ArticlesIds,
		})
	}

	_, err = d.riskFactorsCollection.ReplaceOne(ctx, bson.M{"_id": riskFactorsModel.ID}, riskFactorsModel, options.Replace().SetUpsert(true))
	if err != nil {
		return nil, utils.WrapErrorf(err, "error updating risk factors for disease [%s] in DB", disease.ID)
	}

	return &model.UpdateRiskFactorsPayload{}, nil
}

func convertDiseaseToGql(disease *Disease) *model.Disease {
	gqlDisease := &model.Disease{
		ID:          disease.ID,
		Names:       disease.Names,
		Category:    disease.Category,
		Description: disease.Description,
	}

	if disease.DbLinks != nil {
		gqlDisease.DbLinks = &model.DiseaseDBLinks{
			Icd10: disease.DbLinks.Icd10,
			Icd11: disease.DbLinks.Icd11,
			Mesh:  disease.DbLinks.Mesh,
		}
	}

	return gqlDisease
}

func convertDiseaseToDB(gqlDisease *model.Disease) *Disease {
	disease := &Disease{
		ID:          gqlDisease.ID,
		Names:       gqlDisease.Names,
		Category:    gqlDisease.Category,
		Description: gqlDisease.Description,
	}

	if gqlDisease.DbLinks != nil {
		disease.DbLinks = &DiseaseDBLinks{
			Icd10: gqlDisease.DbLinks.Icd10,
			Icd11: gqlDisease.DbLinks.Icd11,
			Mesh:  gqlDisease.DbLinks.Mesh,
		}
	}

	return disease
}

func convertArticleToGql(article *Article) *model.Article {
	gqlArticle := &model.Article{
		ID:   article.ID,
		Text: article.Text,
	}

	return gqlArticle
}

func convertQAToGql(qa *QA) *model.Qa {
	gqlQA := &model.Qa{
		ID:      qa.ID,
		Disease: convertDiseaseToGql(qa.Disease),
		Article: convertArticleToGql(qa.Article),
	}

	for _, question := range qa.Questions {
		gqlQA.Questions = append(gqlQA.Questions, convertQuestionToGql(question))
	}

	return gqlQA
}

func convertQuestionToGql(question *Question) *model.Question {
	gqlQuestion := &model.Question{
		ID:   question.ID,
		Text: question.Text,
	}

	for _, answer := range question.Answers {
		gqlQuestion.Answers = append(gqlQuestion.Answers, convertAnswerToGql(answer))
	}

	return gqlQuestion
}

func convertAnswerToGql(answer *Answer) *model.Answer {
	gqlAnswer := &model.Answer{
		AnswerStart: answer.AnswerStart,
		Text:        answer.Text,
	}

	return gqlAnswer
}

func convertClassificationItemToGql(classificationItem *ClassificationItem) *model.ClassificationItem {
	gqlClassificationItem := &model.ClassificationItem{
		ID:      classificationItem.ID,
		Article: convertArticleToGql(classificationItem.Article),
		Label:   int(classificationItem.Label),
	}

	return gqlClassificationItem
}

func convertRiskFactorsToGql(riskFactors *RiskFactors) *model.RiskFactors {
	gqlRiskFactors := &model.RiskFactors{
		ID:      riskFactors.ID,
		Disease: convertDiseaseToGql(riskFactors.Disease),
	}

	for _, riskFactor := range riskFactors.RiskFactors {
		gqlRiskFactors.RiskFactors = append(gqlRiskFactors.RiskFactors, &model.RiskFactor{
			Text:        riskFactor.Text,
			Score:       riskFactor.Score,
			ArticlesIds: riskFactor.ArticlesIDs,
		})
	}

	return gqlRiskFactors
}

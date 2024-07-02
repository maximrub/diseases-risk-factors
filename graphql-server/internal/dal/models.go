package dal

type Answer struct {
	AnswerStart int    `bson:"answer_start"`
	Text        string `bson:"text,omitempty"`
}

type Article struct {
	ID   string `bson:"_id,omitempty"`
	Text string `bson:"text,omitempty"`
}

type Disease struct {
	ID          string          `bson:"_id,omitempty"`
	Names       []string        `bson:"names,omitempty"`
	DbLinks     *DiseaseDBLinks `bson:"dbLinks,omitempty"`
	Category    string          `bson:"category,omitempty"`
	Description string          `bson:"description,omitempty"`
}

type DiseaseDBLinks struct {
	Icd10 []string `bson:"icd10,omitempty"`
	Icd11 []string `bson:"icd11,omitempty"`
	Mesh  []string `bson:"mesh,omitempty"`
}

type QA struct {
	ID        string      `bson:"_id,omitempty"`
	Disease   *Disease    `bson:"disease,omitempty"`
	Article   *Article    `bson:"article,omitempty"`
	Questions []*Question `bson:"questions,omitempty"`
}

type Question struct {
	ID      string    `bson:"_id,omitempty"`
	Text    string    `bson:"text,omitempty"`
	Answers []*Answer `bson:"answers,omitempty"`
}

type ClassificationLabel int

const (
	ClassificationLabelNegative ClassificationLabel = 0
	ClassificationLabelPositive ClassificationLabel = 1
)

type ClassificationItem struct {
	ID      string              `bson:"_id,omitempty"`
	Label   ClassificationLabel `bson:"label"`
	Article *Article            `bson:"article,omitempty"`
}

type RiskFactors struct {
	ID          string        `bson:"_id,omitempty"`
	Disease     *Disease      `bson:"disease,omitempty"`
	RiskFactors []*RiskFactor `bson:"riskFactors,omitempty"`
}

type RiskFactor struct {
	Text        string   `bson:"text,omitempty"`
	Score       float64  `bson:"score,omitempty"`
	ArticlesIDs []string `bson:"articlesIDs"`
}

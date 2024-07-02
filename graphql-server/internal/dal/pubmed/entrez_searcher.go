package pubmed

import (
	"encoding/xml"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/biogo/ncbi/entrez"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/utils"
	"io"
	"os"
	"strings"
	"time"
)

const (
	tool = "diseases-risk-factors"
)

type entrezSearcher struct {
	email string
}

func NewEntrezSearcher() *entrezSearcher {
	return &entrezSearcher{
		email: os.Getenv("ENTREZ_EMAIL"),
	}
}

func (s *entrezSearcher) Search(query string, limit int) ([]int, error) {
	const db = "pubmed"
	history := entrez.History{}
	params := &entrez.Parameters{
		RetMax: limit,
	}

	var results *entrez.Search
	err := retry.Do(
		func() error {
			var err error
			results, err = entrez.DoSearch(db, query, params, &history, tool, s.email)
			if err != nil {
				return utils.WrapErrorf(err, "error searching with query \"%s\"", query)
			}

			return nil
		},
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Second),
	)
	if err != nil {
		return nil, utils.WrapErrorf(err, "error searching with query \"%s\" with retry", query)
	}

	if results != nil {
		return results.IdList, nil
	}

	return nil, nil
}

func (s *entrezSearcher) Fetch(id int) (string, error) {
	const db = "pubmed"
	params := &entrez.Parameters{
		RetMode: "xml",
		RetType: "abstract",
	}

	var reader io.ReadCloser
	err := retry.Do(
		func() error {
			var err error
			reader, err = entrez.Fetch(db, params, tool, s.email, nil, id)
			if err != nil {
				return utils.WrapErrorf(err, "error fetch article with id [%s]", id)
			}

			return nil
		},
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Second),
	)
	if err != nil {
		return "", utils.WrapErrorf(err, "error fetch article with id [%s] with retry", id)
	}

	defer reader.Close()
	buf, err := io.ReadAll(reader)
	if err != nil {
		return "", utils.WrapErrorf(err, "error reading article with id [%s]", id)
	}

	var summery ArticleSummery
	err = xml.Unmarshal(buf, &summery)
	if err != nil {
		return "", utils.WrapErrorf(err, "error unmarshal article with id [%s]", id)
	}

	var articleAbstract string
	for _, abstractText := range summery.PubmedArticle.MedlineCitation.Article.Abstract.AbstractText {
		if abstractText.Label != "" {
			articleAbstract += fmt.Sprintf("%s:\n", abstractText.Label)
		}

		text := strings.Join(strings.Fields(strings.TrimSpace(abstractText.Text)), " ")
		articleAbstract += fmt.Sprintf("%s\n\n", text)
	}

	articleAbstract = strings.TrimSuffix(articleAbstract, "\n\n")
	return articleAbstract, nil
}

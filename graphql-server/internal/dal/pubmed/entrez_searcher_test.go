package pubmed_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal/pubmed"
)

type entrezSearcherTestSuite struct {
	suite.Suite
	articleSearcher *pubmed.EntrezSearcher
}

func TestEntrezSearcher(t *testing.T) {
	suite.Run(t, new(entrezSearcherTestSuite))
}

func (t *entrezSearcherTestSuite) SetupTest() {
	t.articleSearcher = pubmed.NewEntrezSearcher()
}

func (t *entrezSearcherTestSuite) TestSanity() {
	r := t.Require()

	ids, err := t.articleSearcher.Search("asthma[mesh]", 1)
	r.NoError(err)
	r.NotEmpty(ids)
	r.Len(ids, 1)

	text, err := t.articleSearcher.Fetch(ids[0])
	r.NoError(err)
	r.NotEmpty(text)
}

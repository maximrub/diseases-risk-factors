package kegg_test

import (
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/dal/kegg"
	"github.com/stretchr/testify/suite"
	"testing"
)

type keggDiseasesFetcherTestSuite struct {
	suite.Suite
	diseasesFetcher *kegg.DiseasesFetcher
}

func TestKeggDiseasesFetcher(t *testing.T) {
	suite.Run(t, new(keggDiseasesFetcherTestSuite))
}

func (t *keggDiseasesFetcherTestSuite) SetupTest() {
	t.diseasesFetcher = kegg.NewDiseasesFetcher()
}

func (t *keggDiseasesFetcherTestSuite) TestSanity() {
	r := t.Require()

	ids, err := t.diseasesFetcher.List()
	r.NoError(err)
	r.NotEmpty(ids)

	disease, err := t.diseasesFetcher.Fetch(ids[0])
	r.NoError(err)
	r.NotNil(disease)
}

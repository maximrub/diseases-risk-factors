package kegg

import (
	"bufio"
	"fmt"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/graph/model"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/utils"
	"net/http"
	"strings"
)

const (
	listUrl = "http://rest.kegg.jp/list/disease"
	getUrl  = "http://rest.kegg.jp/get"
)

type diseasesFetcher struct {
}

func NewDiseasesFetcher() *diseasesFetcher {
	return &diseasesFetcher{}
}

func (f *diseasesFetcher) List() ([]string, error) {
	resp, err := http.Get(listUrl)
	if err != nil {
		return nil, utils.WrapError(err, "error getting diseases IDs")
	}

	defer resp.Body.Close()
	var ids []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		ids = append(ids, strings.Fields(scanner.Text())[0])
	}

	return ids, nil
}

func (f *diseasesFetcher) Fetch(id string) (*model.Disease, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", getUrl, id))
	if err != nil {
		return nil, utils.WrapError(err, "error getting diseases IDs")
	}

	defer resp.Body.Close()
	var disease KeggDisease
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ENTRY") {
			disease.ID = strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, "ENTRY")))[0]
		} else if strings.HasPrefix(line, "NAME") {
			name := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "NAME")), ";")
			disease.Names = append(disease.Names, name)
			for scanner.Scan() {
				subLine := scanner.Text()
				if !strings.HasPrefix(subLine, " ") || strings.HasPrefix(subLine, "  SUBGROUP") {
					break
				}

				subLine = strings.TrimSuffix(strings.TrimSpace(subLine), ";")
				disease.Names = append(disease.Names, subLine)
			}
		} else if strings.HasPrefix(line, "DBLINKS") {
			disease.DBLinks = make(map[string][]string)
			for scanner.Scan() {
				subLine := scanner.Text()
				if !strings.HasPrefix(subLine, " ") {
					break
				}
				subLine = strings.TrimSpace(subLine)
				parts := strings.SplitN(subLine, ":", 2)
				disease.DBLinks[parts[0]] = strings.Fields(parts[1])
			}
		} else if strings.HasPrefix(line, "CATEGORY") {
			disease.Category = strings.TrimSpace(strings.TrimPrefix(line, "CATEGORY"))
		} else if strings.HasPrefix(line, "DESCRIPTION") {
			disease.Description = strings.TrimSpace(strings.TrimPrefix(line, "DESCRIPTION "))
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, utils.WrapErrorf(err, "error reading disease [%s]", id)
	}

	return &model.Disease{
		ID:    disease.ID,
		Names: disease.Names,
		DbLinks: &model.DiseaseDBLinks{
			Icd10: disease.DBLinks["ICD-10"],
			Icd11: disease.DBLinks["ICD-11"],
			Mesh:  disease.DBLinks["MeSH"],
		},
		Category:    disease.Category,
		Description: disease.Description,
	}, nil
}

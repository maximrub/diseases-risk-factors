package kegg

type KeggDisease struct {
	ID          string
	Names       []string
	DBLinks     map[string][]string
	Category    string
	Description string
}

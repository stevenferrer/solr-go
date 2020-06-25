package query

// Facet is a facet
// https://lucene.apache.org/solr/guide/8_5/json-facet-api.html
type Facet struct {
	Field  string
	Type   string
	Limit  int
	Sort   string
	Facet  map[string]interface{}
	Domain *FDomain
}

// FDomain is a facet Domain
type FDomain struct {
	ExcludeTags string
	Filters     []string
}

// NewFacet is a factory for *Facet
func NewFacet(field, typeRes string) *Facet {
	return &Facet{Field: field, Type: typeRes}
}

func (f *Facet) AddNestedFacet(label string, nf *Facet) {
	f.Facet[label] = nf
}

func (f *Facet) SetLimit(limit int) {
	f.Limit = limit
}

func (f *Facet) SetSort(sort string) {
	f.Sort = sort
}

func (f *Facet) SetDomain(excludeTags string, filters []string) {
	f.Domain = &FDomain{
		ExcludeTags: excludeTags,
		Filters:     filters,
	}
}

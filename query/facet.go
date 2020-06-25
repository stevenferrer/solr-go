package query

// Facet is a facet
type Facet struct {
	Field  string
	Type   string
	Facet  map[string]interface{}
	Domain *FDomain
}

// FDomain is a facet Domain
type FDomain struct {
	ExcludeTags string
	Filters     []string
}

func NewFacet(field, typeRes string) *Facet {
	return &Facet{Field: field, Type: typeRes}
}

func (f *Facet) AddFacet(k, v string) {
	f.Facet[k] = v
}

func (f *Facet) SetDomain(excludeTags string, filters []string) {
	f.Domain = &FDomain{
		ExcludeTags: excludeTags,
		Filters:     filters,
	}
}

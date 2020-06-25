package query

// Facet is a facet
type Facet struct {
	Field  string
	Type   string
	Facet  map[string]interface{}
	Domain struct {
		ExcludeTags string
		Filters     []string
	}
}

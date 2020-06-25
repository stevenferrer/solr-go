package query

// Query is a query
type Query struct {
	facets   map[string]Facet
	queryStr string
	queries  map[string]interface{}
	filters  []string
	fields   []string
}

func NewQuery() *Query {
	return &Query{}
}

// SetQuery sets the main query
func (q *Query) SetQuery(qs string) {}

// AddQuery adds a query in `queries`
func (q *Query) AddQuery(k string, v interface{}) {}

// AddFilter adds a filter
func (q *Query) AddFilter(filter string) {}

// AddField adds a field
func (q *Query) AddField(field string) {}

// AddFacet adds a facet
func (q *Query) AddFacet(facet Facet) {}

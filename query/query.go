package query

// Query is a query
type Query struct {
	qs      string
	queries map[string]interface{}
	filters []string
	fields  []string
	facets  map[string]*Facet
}

// NewQuery is a factory for Query
func NewQuery(qs string) *Query {
	return &Query{qs: qs}
}

// SetQuery sets the main query
func (q *Query) SetQuery(qs string) {
	q.qs = qs
}

// AddQuery adds a query in `queries`
func (q *Query) AddQuery(k string, v interface{}) {
	q.queries[k] = v
}

// AddFilter adds a filter
func (q *Query) AddFilter(filter string) {
	q.filters = append(q.filters, filter)
}

// AddField adds a field
func (q *Query) AddField(field string) {
	q.fields = append(q.fields, field)
}

// AddFacet adds a facet
func (q *Query) AddFacet(label string, facet *Facet) {
	q.facets[label] = facet
}

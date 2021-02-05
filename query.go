package solr

// Query is a query
type Query struct {
	// common query params
	// Refer to https://lucene.apache.org/solr/guide/8_7/common-query-parameters.html

	// standard params in v1 api
	// debug                 string
	// explainOther          string
	// timeAllowed           int
	// segmentTerminateEarly bool
	// omitHeader            bool
	// echoParams            string
	// minExactCount         int
	// responseWriter        string // wt

	// supported params in json request api
	sort    string
	offset  int      // start
	limit   int      // rows
	filters []string // fq
	fields  []string // fl

	// query parser
	// Refer to https://lucene.apache.org/solr/guide/8_7/query-syntax-and-parsing.html
	qp QueryParser

	// facets
	// Refer to https://lucene.apache.org/solr/guide/8_7/json-facet-api.html
	facets []Faceter

	// additional queries
	// https://lucene.apache.org/solr/guide/8_7/json-query-dsl.html#additional-queries
	queries M
}

// NewQuery returns a new Query
func NewQuery() *Query {
	return &Query{}
}

// BuildQuery builds the query
func (q *Query) BuildQuery() M {
	qm := M{"query": q.qp.BuildParser()}

	if q.queries != nil {
		qm["queries"] = q.queries
	}

	if q.sort != "" {
		qm["sort"] = q.sort
	}

	if q.offset != 0 {
		qm["offset"] = q.offset
	}

	if q.limit != 0 {
		qm["limit"] = q.limit
	}

	if len(q.filters) > 0 {
		qm["filter"] = q.filters
	}

	if len(q.fields) != 0 {
		qm["fields"] = q.fields
	}

	if len(q.facets) > 0 {
		facets := M{}
		for _, facet := range q.facets {
			facets[facet.Name()] = facet.BuildFacet()
		}

		qm["facet"] = facets
	}

	return qm
}

// Sort sets the sort param
func (q *Query) Sort(sort string) *Query {
	q.sort = sort
	return q
}

// Offset sets the offset param
func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

// Limit sets the limit param
func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

// Filters sets the filter param
func (q *Query) Filters(filters ...string) *Query {
	q.filters = filters
	return q
}

// Fields sets the fields param
func (q *Query) Fields(fields ...string) *Query {
	q.fields = fields
	return q
}

// QueryParser sets the query parser
func (q *Query) QueryParser(qp QueryParser) *Query {
	q.qp = qp
	return q
}

// Facets sets the facet query
func (q *Query) Facets(facets ...Faceter) *Query {
	q.facets = facets
	return q
}

// Queries sets the additional queries
func (q *Query) Queries(queries M) *Query {
	q.queries = queries
	return q
}

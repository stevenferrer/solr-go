package solr

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Query is a query
type Query struct {
	// common query params
	// https://lucene.apache.org/solr/guide/8_7/common-query-parameters.html

	// params in standard api
	debug                 string
	explainOther          string
	timeAllowed           int
	segmentTerminateEarly bool
	omitHeader            bool
	echoParams            string
	minExactCount         int
	responseWriter        string // wt

	// supported params in json request api
	sort   string
	offset int    // start
	limit  int    // rows
	filter string // fq
	fields string // fl

	// query parser
	// https://lucene.apache.org/solr/guide/8_7/query-syntax-and-parsing.html
	qp QueryParser

	// facet query
	// https://lucene.apache.org/solr/guide/8_7/faceting.html
	// https://lucene.apache.org/solr/guide/8_7/json-facet-api.html
	facets Facets

	// additional queries
	// https://lucene.apache.org/solr/guide/8_7/json-query-dsl.html#additional-queries
	queries M
}

// NewQuery returns a new Query
func NewQuery() *Query {
	return &Query{}
}

// BuildJSON builds the query to JSON
func (q *Query) BuildJSON() ([]byte, error) {
	qq, err := q.qp.BuildParser()
	if err != nil {
		return nil, errors.Wrap(err, "build parser")
	}

	qm := M{"query": qq}

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

	if q.filter != "" {
		qm["filter"] = q.filter
	}

	if q.fields != "" {
		qm["fields"] = q.fields
	}

	b, err := json.Marshal(qm)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal query")
	}

	return b, nil
}

// WithSort sets the sort param
func (q *Query) WithSort(sort string) *Query {
	q.sort = sort
	return q
}

// WithOffset sets the offset param
func (q *Query) WithOffset(offset int) *Query {
	q.offset = offset
	return q
}

// WithLimit sets the limit param
func (q *Query) WithLimit(limit int) *Query {
	q.limit = limit
	return q
}

// WithFilter sets the filter param
func (q *Query) WithFilter(filter string) *Query {
	q.filter = filter
	return q
}

// WithFields sets the fields param
func (q *Query) WithFields(fields string) *Query {
	q.fields = fields
	return q
}

// WithQueryParser sets the query parser
func (q *Query) WithQueryParser(qp QueryParser) *Query {
	q.qp = qp
	return q
}

// WithFacets sets the facet query
func (q *Query) WithFacets(facets Facets) *Query {
	q.facets = facets
	return q
}

// WithQueries sets the additional queries
func (q *Query) WithQueries(queries M) *Query {
	q.queries = queries
	return q
}

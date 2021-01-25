package solr

import (
	"github.com/pkg/errors"
)

// QueryBuilder is a query builder
type QueryBuilder struct {
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

// NewQueryBuilder returns a new QueryBuilder
func NewQueryBuilder(qp QueryParser) *QueryBuilder {
	return &QueryBuilder{qp: qp}
}

// Build builds the query
func (qb *QueryBuilder) Build() (M, error) {
	qq, err := qb.qp.BuildQuery()
	if err != nil {
		return nil, errors.Wrap(err, "build query from parser")
	}

	return M{"query": qq}, nil
}

// WithSort sets the sort param
func (qb *QueryBuilder) WithSort(sort string) *QueryBuilder {
	qb.sort = sort
	return qb
}

// WithOffset sets the offset param
func (qb *QueryBuilder) WithOffset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// WithLimit sets the limit param
func (qb *QueryBuilder) WithLimit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// WithFilter sets the filter param
func (qb *QueryBuilder) WithFilter(filter string) *QueryBuilder {
	qb.filter = filter
	return qb
}

// WithFields sets the fields param
func (qb *QueryBuilder) WithFields(fields string) *QueryBuilder {
	qb.fields = fields
	return qb
}

// WithQueryParser sets the query parser
func (qb *QueryBuilder) WithQueryParser(qp QueryParser) *QueryBuilder {
	qb.qp = qp
	return qb
}

// WithFacets sets the facet query
func (qb *QueryBuilder) WithFacets(facets Facets) *QueryBuilder {
	qb.facets = facets
	return qb
}

// WithAdditionalQueries sets the additional queries
func (qb *QueryBuilder) WithAdditionalQueries(queries M) *QueryBuilder {
	qb.queries = queries
	return qb
}

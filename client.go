package solr

import (
	"context"
	"errors"
)

// Client is an abstraction of a solr client
// e.g. standard, json api or solr cloud (v2 api)
type Client interface {
	Query(ctx context.Context, collection string, query *QueryBuilder) (*QueryResponse, error)
}

// JSONClient is a solr client that uses the JSON API
type JSONClient struct{}

var _ Client = (*JSONClient)(nil)

// Query is used for searching/querying
func (c *JSONClient) Query(ctx context.Context, collection string, qb *QueryBuilder) (*QueryResponse, error) {
	return nil, errors.New("not implemented")
}

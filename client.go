package solr

import (
	"context"
)

// Client is an abstraction of a solr client e.g. standard, json api or solr cloud (v2 api)
type Client interface {
	Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error)
	Update(ctx context.Context, collection string, documents ...Document) (*UpdateResponse, error)
	Commit(ctx context.Context, collection string) error
}

package solr

import (
	"context"
	"io"
)

// Client is an abstraction of a solr client e.g. standard, json api or v2 api
type Client interface {
	Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error)

	Update(ctx context.Context, collection string, ct ContentType, body io.Reader) (*UpdateResponse, error)
	Commit(ctx context.Context, collection string) error

	AddFields(ctx context.Context, collection string, fields ...Field) error
	DeleteFields(ctx context.Context, collection string, fields ...Field) error
	ReplaceFields(ctx context.Context, collection string, fields ...Field) error
	AddDynamicFields(ctx context.Context, collection string, fields ...Field) error
	DeleteDynamicFields(ctx context.Context, collection string, fields ...Field) error
	ReplaceDynamicFields(ctx context.Context, collection string, fields ...Field) error
	AddFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	DeleteFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	ReplaceFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	AddCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error
	DeleteCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error

	SetProperties(ctx context.Context, collection string, properties ...CommonProperty) error
	UnsetProperty(ctx context.Context, collection string, property CommonProperty) error
	AddComponents(ctx context.Context, collection string, component ...*Component) error
	UpdateComponents(ctx context.Context, collection string, component ...*Component) error
	DeleteComponents(ctx context.Context, collection string, component ...*Component) error

	Suggest(ctx context.Context, collection string, params *SuggestParams) (*SuggestResponse, error)
}

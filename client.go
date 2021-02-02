package solr

import (
	"context"
)

// Client is an abstraction of a solr client e.g. standard, json api or solr cloud (v2 api)
type Client interface {
	Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error)

	Update(ctx context.Context, collection string, documents ...Document) (*UpdateResponse, error)
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

	AddComponent(ctx context.Context, collection string, component Component) error
	UpdateComponent(ctx context.Context, collection string, component Component) error
	DeleteComponent(ctx context.Context, collection string, component Component) error
}

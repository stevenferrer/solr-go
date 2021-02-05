package solr

import (
	"context"
	"io"
)

// Client is an a client interface for interacting with Solr
type Client interface {
	// CreateCollection creates a new collection.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/collection-management.html#create
	CreateCollection(ctx context.Context, params *CollectionParams) error
	// DeleteCollection deletes a collection.
	// Refer to https://lucene.apache.org/solr/guide/8_8/collection-management.html#delete
	DeleteCollection(ctx context.Context, params *CollectionParams) error

	// Query sends a query request to query API.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/json-request-api.html
	Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error)

	// Update can be used to add, update, or delete a document from the index.
	// `body` is expected to contain the list of documents.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/uploading-data-with-index-handlers.html
	Update(ctx context.Context, collection string, ct ContentType, body io.Reader) (*UpdateResponse, error)
	// Commit commits the last update
	Commit(ctx context.Context, collection string) error

	// AddFields adds new field definitions to the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-field
	AddFields(ctx context.Context, collection string, fields ...Field) error
	// DeleteFields removes field definitions from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-field
	DeleteFields(ctx context.Context, collection string, fields ...Field) error
	// ReplaceFields replaces field definition from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-field
	ReplaceFields(ctx context.Context, collection string, fields ...Field) error
	// AddDynamicFields adds new dynamic field rules to the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-dynamic-field-rule
	AddDynamicFields(ctx context.Context, collection string, fields ...Field) error
	// DeleteDynamicFields removes dynamic field rules from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-dynamic-field-rule
	DeleteDynamicFields(ctx context.Context, collection string, fields ...Field) error
	// ReplaceDynamicFields replaces dynamic field rules from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-dynamic-field-rule
	ReplaceDynamicFields(ctx context.Context, collection string, fields ...Field) error
	// AddFieldTypes adds new field types to the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-field-type
	AddFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	// DeleteFieldTypes removes field type definitions from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-field-type
	DeleteFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	// ReplaceFieldTypes replaces field type defintions from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-field-type
	ReplaceFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error
	// AddCopyFields adds new copy field rules to the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-copy-field-rule
	AddCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error
	// DeleteCopyFields deletes copy field rules from the schema.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-copy-field-rule
	DeleteCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error

	// SetProperties sets well known properties.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-common-properties
	SetProperties(ctx context.Context, collection string, properties ...CommonProperty) error
	// UnsetProperty removes a property set via SetProperties.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-common-properties
	UnsetProperty(ctx context.Context, collection string, property CommonProperty) error
	// AddComponents adds a component (request handler, search component, init params, etc.) to configoverlay.json.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
	AddComponents(ctx context.Context, collection string, component ...*Component) error
	// UpdateComponents overwrites existing settings from configoverlay.json.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
	UpdateComponents(ctx context.Context, collection string, component ...*Component) error
	// DeleteComponents removes settings from configoverlay.json
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
	DeleteComponents(ctx context.Context, collection string, component ...*Component) error

	// Suggest queries the suggest endpoint.
	//
	// Refer to https://lucene.apache.org/solr/guide/8_8/suggester.html#get-suggestions-with-weights
	Suggest(ctx context.Context, collection string, params *SuggestParams) (*SuggestResponse, error)
}

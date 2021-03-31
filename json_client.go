package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// JSONClient is a client for interacting with Solr via JSON API
type JSONClient struct {
	// baseURL is the base url of the solr instance
	baseURL string
	// reqSender is the request sender
	reqSender RequestSender
}

var _ Client = (*JSONClient)(nil)

// NewJSONClient returns a new JSONClient
func NewJSONClient(baseURL string) *JSONClient {
	return &JSONClient{
		baseURL:   baseURL,
		reqSender: NewDefaultRequestSender(),
	}
}

// WithRequestSender overrides the default request sender
func (c *JSONClient) WithRequestSender(reqSender RequestSender) *JSONClient {
	c.reqSender = reqSender
	return c
}

// CreateCollection creates a new collection.
//
// Refer to https://solr.apache.org/guide/8_8/collection-management.html#create
func (c *JSONClient) CreateCollection(ctx context.Context, params *CollectionParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/collections?action=CREATE&"+params.BuildParams(), c.baseURL)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// DeleteCollection deletes a collection.
//
// Refer to https://solr.apache.org/guide/8_8/collection-management.html#delete
func (c *JSONClient) DeleteCollection(ctx context.Context, params *CollectionParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/collections?action=DELETE&"+params.BuildParams(), c.baseURL)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// CoreStatus returns the status of all running Solr cores, or status for only the named core.
//
// Refer to https://solr.apache.org/guide/8_8/coreadmin-api.html#coreadmin-status
func (c *JSONClient) CoreStatus(ctx context.Context, params *CoreParams) (*CoreStatusResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/admin/cores?action=STATUS&"+params.BuildParams(), c.baseURL)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp CoreStatusResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	return &resp, nil
}

// CreateCore creates a new core
//
// Refer to https://solr.apache.org/guide/8_8/coreadmin-api.html#coreadmin-create
func (c *JSONClient) CreateCore(ctx context.Context, params *CreateCoreParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/cores?action=CREATE&"+params.BuildParams(), c.baseURL)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// UnloadCore removes a core from Solr
//
// Refer to https://solr.apache.org/guide/8_8/coreadmin-api.html#coreadmin-unload
func (c *JSONClient) UnloadCore(ctx context.Context, params *CoreParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/cores?action=UNLOAD&"+params.BuildParams(), c.baseURL)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// Query sends a query request to query API.
//
// Refer to https://solr.apache.org/guide/8_8/json-request-api.html
func (c *JSONClient) Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(query.BuildQuery())
	if err != nil {
		return nil, errors.Wrap(err, "encode request body")
	}

	urlStr := fmt.Sprintf("%s/solr/%s/query", c.baseURL, collection)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodPost, urlStr, JSON.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp QueryResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	return &resp, nil
}

// Update can be used to add, update, or delete a document from the index.
//
// Refer to https://solr.apache.org/guide/8_8/uploading-data-with-index-handlers.html
func (c *JSONClient) Update(ctx context.Context, collection string, mimeType MimeType, body io.Reader) (*UpdateResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/update", c.baseURL, collection)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodPost, urlStr, mimeType.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp UpdateResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	return &resp, nil
}

// Commit commits the last update.
func (c *JSONClient) Commit(ctx context.Context, collection string) error {
	urlStr := fmt.Sprintf("%s/solr/%s/update?commit=true", c.baseURL, collection)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp UpdateResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// AddFields adds new field definitions to the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#add-a-new-field
func (c *JSONClient) AddFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-field", fields)
}

// DeleteFields removes field definitions from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#delete-a-field
func (c *JSONClient) DeleteFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-field", fields)
}

// ReplaceFields replaces field definition from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#replace-a-field
func (c *JSONClient) ReplaceFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-field", fields)
}

// AddDynamicFields adds new dynamic field rules to the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#add-a-dynamic-field-rule
func (c *JSONClient) AddDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-dynamic-field", fields)
}

// DeleteDynamicFields removes dynamic field rules from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#delete-a-dynamic-field-rule
func (c *JSONClient) DeleteDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-dynamic-field", fields)
}

// ReplaceDynamicFields replaces dynamic field rules from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#replace-a-dynamic-field-rule
func (c *JSONClient) ReplaceDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-dynamic-field", fields)
}

// AddFieldTypes adds new field types to the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#add-a-new-field-type
func (c *JSONClient) AddFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "add-field-type", fieldTypes)
}

// DeleteFieldTypes removes field type definitions from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#delete-a-field-type
func (c *JSONClient) DeleteFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "delete-field-type", fieldTypes)
}

// ReplaceFieldTypes replaces field type defintions from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#replace-a-field-type
func (c *JSONClient) ReplaceFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "replace-field-type", fieldTypes)
}

// AddCopyFields adds new copy field rules to the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#add-a-new-copy-field-rule
func (c *JSONClient) AddCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "add-copy-field", copyFields)
}

// DeleteCopyFields deletes copy field rules from the schema.
//
// Refer to https://solr.apache.org/guide/8_8/schema-api.html#delete-a-copy-field-rule
func (c *JSONClient) DeleteCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "delete-copy-field", copyFields)
}

func (c *JSONClient) modifySchema(ctx context.Context, collection, command string, body interface{}) error {
	urlStr := fmt.Sprintf("%s/solr/%s/schema", c.baseURL, collection)
	return c.postJSON(ctx, urlStr, M{command: body})
}

func (c *JSONClient) postJSON(ctx context.Context, urlStr string, reqBody interface{}) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		return errors.Wrap(err, "encode request body")
	}

	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodPost, urlStr, JSON.String(), buf)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil
}

// SetProperties sets well known properties.
//
// Refer to https://solr.apache.org/guide/8_8/config-api.html#commands-for-common-properties
func (c *JSONClient) SetProperties(ctx context.Context, collection string, properties ...CommonProperty) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	m := M{}
	for _, prop := range properties {
		m[prop.Name] = prop.Value
	}

	return c.postJSON(ctx, urlStr, M{"set-property": m})
}

// UnsetProperty removes a property set via set-properties.
//
// Refer to https://solr.apache.org/guide/8_8/config-api.html#commands-for-common-properties
func (c *JSONClient) UnsetProperty(ctx context.Context, collection string, property CommonProperty) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	return c.postJSON(ctx, urlStr, M{"unset-property": property.Name})
}

// AddComponents adds a component (request handler, search component, init params, etc.) to configoverlay.json.
//
// Refer to https://solr.apache.org/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) AddComponents(ctx context.Context, collection string, components ...*Component) error {
	commands := []string{}
	for _, comp := range components {
		b, err := json.Marshal(comp.BuildComponent())
		if err != nil {
			return errors.Wrap(err, "marshal component")
		}

		command := fmt.Sprintf("%q:%s", "add-"+comp.ct.String(), string(b))
		commands = append(commands, command)
	}

	reqBody := "{" + strings.Join(commands, ",") + "}"

	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodPost, urlStr, JSON.String(), strings.NewReader(reqBody))
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	return nil

}

// UpdateComponents overwrites existing settings from configoverlay.json.
//
// Refer to https://solr.apache.org/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) UpdateComponents(ctx context.Context, collection string, components ...*Component) error {
	return errors.New("not implemented")
}

// DeleteComponents removes settings from configoverlay.json
//
// Refer to https://solr.apache.org/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) DeleteComponents(ctx context.Context, collection string, components ...*Component) error {
	return errors.New("not implemented")
}

// Suggest queries the suggest endpoint.
//
// Refer to https://solr.apache.org/guide/8_8/suggester.html#get-suggestions-with-weights
func (c *JSONClient) Suggest(ctx context.Context, collection string, params *SuggestParams) (*SuggestResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/%s?%s", c.baseURL, collection, params.endpoint, params.BuildParams())
	httpResp, err := c.reqSender.SendRequest(ctx, http.MethodGet, urlStr, JSON.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp SuggestResponse
	err = readResponse(httpResp, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	return &resp, nil
}

func readResponse(resp *http.Response, v interface{}) error {
	contentType := resp.Header.Get("content-type")
	if strings.Contains(contentType, "text/html") {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read html response")
		}

		return errors.Errorf("unexpected html response: %s", string(b))
	}

	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return errors.Wrap(err, "decode json response")
	}

	val, ok := v.(*BaseResponse)
	if ok && resp.StatusCode > http.StatusOK {
		return val.Error
	}

	return nil
}

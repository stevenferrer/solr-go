package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// JSONClient is a client for interacting with Apache Solr via JSON API
type JSONClient struct {
	// baseURL is the base url of the solr instance
	baseURL    string
	httpClient *http.Client
}

var _ Client = (*JSONClient)(nil)

// NewJSONClient returns a new JSONClient
func NewJSONClient(baseURL string) *JSONClient {
	return &JSONClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

// WithHTTPClient overrides the default http client
func (c *JSONClient) WithHTTPClient(httpClient *http.Client) *JSONClient {
	c.httpClient = httpClient
	return c
}

// CreateCollection creates a new collection.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/collection-management.html#create
func (c *JSONClient) CreateCollection(ctx context.Context, params *CollectionParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/collections?action=CREATE&"+params.BuildParam(), c.baseURL)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	httpResp, err := c.sendRequest(ctx, http.MethodGet, theURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil
}

// DeleteCollection deletes a collection.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/collection-management.html#delete
func (c *JSONClient) DeleteCollection(ctx context.Context, params *CollectionParams) error {
	urlStr := fmt.Sprintf("%s/solr/admin/collections?action=DELETE&"+params.BuildParam(), c.baseURL)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	httpResp, err := c.sendRequest(ctx, http.MethodGet, theURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil
}

// Query sends a query request to query API.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/json-request-api.html
func (c *JSONClient) Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/query", c.baseURL, collection)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(query.BuildQuery())
	if err != nil {
		return nil, errors.Wrap(err, "encode query")
	}

	httpResp, err := c.sendRequest(ctx, http.MethodPost, theURL.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp QueryResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

// Update can be used to add, update, or delete a document from the index.
// `body` is expected to contain the list of documents.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/uploading-data-with-index-handlers.html
func (c *JSONClient) Update(ctx context.Context, collection string, ct ContentType, body io.Reader) (*UpdateResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/update", c.baseURL, collection)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	httpResp, err := c.sendRequestWithContentType(ctx, http.MethodPost, theURL.String(), ct.String(), body)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var resp UpdateResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

// Commit commits the last update.
func (c *JSONClient) Commit(ctx context.Context, collection string) error {
	urlStr := fmt.Sprintf("%s/solr/%s/update", c.baseURL, collection)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	q := theURL.Query()
	q.Add("commit", "true")
	theURL.RawQuery = q.Encode()

	httpResp, err := c.sendRequest(ctx, http.MethodGet, theURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var res UpdateResponse
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return res.Error
	}

	return nil
}

// AddFields adds new field definitions to the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-field
func (c *JSONClient) AddFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-field", fields)
}

// DeleteFields removes field definitions from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-field
func (c *JSONClient) DeleteFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-field", fields)
}

// ReplaceFields replaces field definition from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-field
func (c *JSONClient) ReplaceFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-field", fields)
}

// AddDynamicFields adds new dynamic field rules to the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-dynamic-field-rule
func (c *JSONClient) AddDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-dynamic-field", fields)
}

// DeleteDynamicFields removes dynamic field rules from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-dynamic-field-rule
func (c *JSONClient) DeleteDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-dynamic-field", fields)
}

// ReplaceDynamicFields replaces dynamic field rules from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-dynamic-field-rule
func (c *JSONClient) ReplaceDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-dynamic-field", fields)
}

// AddFieldTypes adds new field types to the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-field-type
func (c *JSONClient) AddFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "add-field-type", fieldTypes)
}

// DeleteFieldTypes removes field type definitions from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-field-type
func (c *JSONClient) DeleteFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "delete-field-type", fieldTypes)
}

// ReplaceFieldTypes replaces field type defintions from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#replace-a-field-type
func (c *JSONClient) ReplaceFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "replace-field-type", fieldTypes)
}

// AddCopyFields adds new copy field rules to the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#add-a-new-copy-field-rule
func (c *JSONClient) AddCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "add-copy-field", copyFields)
}

// DeleteCopyFields deletes copy field rules from the schema.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/schema-api.html#delete-a-copy-field-rule
func (c *JSONClient) DeleteCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "delete-copy-field", copyFields)
}

func (c *JSONClient) modifySchema(ctx context.Context, collection, command string, body interface{}) error {
	urlStr := fmt.Sprintf("%s/solr/%s/schema", c.baseURL, collection)
	return c.postJSON(ctx, urlStr, M{command: body})
}

func (c *JSONClient) postJSON(
	ctx context.Context,
	urlStr string,
	reqBody interface{},
) error {
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(reqBody)
	if err != nil {
		return errors.Wrap(err, "encode request body")
	}

	httpResp, err := c.sendRequest(ctx, http.MethodPost, theURL.String(), buf)
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var res BaseResponse
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return res.Error
	}

	return nil
}

// SetProperties sets well known properties.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-common-properties
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
// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-common-properties
func (c *JSONClient) UnsetProperty(ctx context.Context, collection string, property CommonProperty) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	return c.postJSON(ctx, urlStr, M{"unset-property": property.Name})
}

// AddComponents adds a component (request handler, search component, init params, etc.) to configoverlay.json.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) AddComponents(ctx context.Context, collection string, components ...*Component) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

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

	httpResp, err := c.sendRequest(ctx, http.MethodPost, theURL.String(), strings.NewReader(reqBody))
	if err != nil {
		return errors.Wrap(err, "send request")
	}

	var resp BaseResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil

}

// UpdateComponents overwrites existing settings from configoverlay.json.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) UpdateComponents(ctx context.Context, collection string, components ...*Component) error {
	return errors.New("not implemented")
}

// DeleteComponents removes settings from configoverlay.json
//
// Refer to https://lucene.apache.org/solr/guide/8_8/config-api.html#commands-for-handlers-and-components
func (c *JSONClient) DeleteComponents(ctx context.Context, collection string, components ...*Component) error {
	return errors.New("not implemented")
}

// Suggest queries the suggest endpoint.
//
// Refer to https://lucene.apache.org/solr/guide/8_8/suggester.html#get-suggestions-with-weights
func (c *JSONClient) Suggest(ctx context.Context, collection string, params *SuggestParams) (*SuggestResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/%s", c.baseURL, collection, params.endpoint)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	theURL.RawQuery = params.BuildParams()

	httpResp, err := c.sendRequest(ctx, http.MethodGet, theURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}

	var res SuggestResponse
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, res.Error
	}

	return &res, nil
}

func (c *JSONClient) sendRequest(ctx context.Context, httpMethod, urlStr string, body io.Reader) (*http.Response, error) {
	return c.sendRequestWithContentType(ctx, httpMethod, urlStr, JSON.String(), body)
}

func (c *JSONClient) sendRequestWithContentType(ctx context.Context, httpMethod, urlStr, contentType string, body io.Reader) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx,
		httpMethod, urlStr, body)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("Content-Type", contentType)

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
	}

	return httpResp, nil
}

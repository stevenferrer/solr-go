package solr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// JSONClient is a client for interacting with Apache Solr using via JSON API
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

// Query is used for querying documents
func (c *JSONClient) Query(
	ctx context.Context,
	collection string,
	q *Query,
) (*QueryResponse, error) {
	urll, err := url.Parse(fmt.Sprintf("%s/solr/%s/query", c.baseURL, collection))
	if err != nil {
		return nil, errors.Wrap(err, "build request url")
	}

	b, err := q.BuildJSON()
	if err != nil {
		return nil, errors.Wrap(err, "build query body")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urll.String(), bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("Content-Type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
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

// Update is used for updating/indexing documents
func (c *JSONClient) Update(
	ctx context.Context,
	collection string,
	documents ...Document,
) (*UpdateResponse, error) {
	urll, err := url.Parse(fmt.Sprintf("%s/solr/%s/update", c.baseURL, collection))
	if err != nil {
		return nil, errors.Wrap(err, "build request url")
	}

	b, err := json.Marshal(documents)
	if err != nil {
		return nil, errors.Wrap(err, "marshal documents")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urll.String(), bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("Content-Type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
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

// Commit commits the last update
func (c *JSONClient) Commit(ctx context.Context, collection string) error {
	urll, err := url.Parse(fmt.Sprintf("%s/solr/%s/update", c.baseURL, collection))
	if err != nil {
		return errors.Wrap(err, "build request url")
	}

	q := urll.Query()
	q.Add("commit", "true")
	urll.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodGet, urll.String(), nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "send http request")
	}

	var resp UpdateResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response body")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil
}

// AddFields adds one or more fields
func (c *JSONClient) AddFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-field", fields)
}

// DeleteFields deletes one or more fields
func (c *JSONClient) DeleteFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-field", fields)
}

// ReplaceFields updates one or more fields
func (c *JSONClient) ReplaceFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-field", fields)
}

// AddDynamicFields adds on ore more dynamic fields
func (c *JSONClient) AddDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "add-dynamic-field", fields)
}

// DeleteDynamicFields deletes one ore more dynamic fields
func (c *JSONClient) DeleteDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "delete-dynamic-field", fields)
}

// ReplaceDynamicFields updates on or more dynamic fields
func (c *JSONClient) ReplaceDynamicFields(ctx context.Context, collection string, fields ...Field) error {
	return c.modifySchema(ctx, collection, "replace-dynamic-field", fields)
}

// AddFieldTypes adds on more more field types
func (c *JSONClient) AddFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "add-field-type", fieldTypes)
}

// DeleteFieldTypes deletes on or more field types
func (c *JSONClient) DeleteFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "delete-field-type", fieldTypes)
}

// ReplaceFieldTypes updates on or more field types
func (c *JSONClient) ReplaceFieldTypes(ctx context.Context, collection string, fieldTypes ...FieldType) error {
	return c.modifySchema(ctx, collection, "replace-field-type", fieldTypes)
}

// AddCopyFields adds on ore more copy fields
func (c *JSONClient) AddCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "add-copy-field", copyFields)
}

// DeleteCopyFields deletes on more more copy fields
func (c *JSONClient) DeleteCopyFields(ctx context.Context, collection string, copyFields ...CopyField) error {
	return c.modifySchema(ctx, collection, "delete-copy-field", copyFields)
}

func (c *JSONClient) modifySchema(ctx context.Context, collection, command string, body interface{}) error {
	urlStr := fmt.Sprintf("%s/solr/%s/schema", c.baseURL, collection)
	return c.commonRequest(ctx, urlStr, M{command: body})
}

func (c *JSONClient) commonRequest(
	ctx context.Context,
	urlStr string,
	reqBody interface{},
) error {
	theURL, err := url.Parse(urlStr)
	if err != nil {
		errors.Wrap(err, "parse request url")
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return errors.Wrap(err, "marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, theURL.String(), bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("Content-Type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "send http request")
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

// SetProperties sets the common config properties
func (c *JSONClient) SetProperties(ctx context.Context, collection string, properties ...CommonProperty) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	m := M{}
	for _, prop := range properties {
		m[prop.Name] = prop.Value
	}

	return c.commonRequest(ctx, urlStr, M{"set-property": m})
}

// UnsetProperty unsets a common config property
func (c *JSONClient) UnsetProperty(ctx context.Context, collection string, property CommonProperty) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	return c.commonRequest(ctx, urlStr, M{"unset-property": property.Name})
}

// AddComponent adds a component
func (c *JSONClient) AddComponent(ctx context.Context, collection string, component Component) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	command := "add-" + component.Type.String()

	m := M{"name": component.Name, "class": component.Class}
	for k, v := range component.Values {
		m[k] = v
	}

	return c.commonRequest(ctx, urlStr, M{command: m})
}

// UpdateComponent updates a component
func (c *JSONClient) UpdateComponent(ctx context.Context, collection string, component Component) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	command := "update-" + component.Type.String()

	m := M{"name": component.Name, "class": component.Class}
	for k, v := range component.Values {
		m[k] = v
	}

	return c.commonRequest(ctx, urlStr, M{command: m})
}

// DeleteComponent deletes a component
func (c *JSONClient) DeleteComponent(ctx context.Context, collection string, component Component) error {
	urlStr := fmt.Sprintf("%s/solr/%s/config", c.baseURL, collection)
	command := "delete-" + component.Type.String()
	return c.commonRequest(ctx, urlStr, M{command: component.Name})
}

// Suggest queries the suggest endpoint
func (c *JSONClient) Suggest(ctx context.Context, collection string, params *SuggestParams) (*SuggestResponse, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s/solr/%s/%s", c.baseURL, collection, params.endpoint))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	theURL.RawQuery = params.BuildParams()

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodGet, theURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
	}

	var resp SuggestResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

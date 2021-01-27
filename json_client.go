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
		return nil, errors.Wrap(err, "do http request")
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
		return nil, errors.Wrap(err, "do http request")
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
		return errors.Wrap(err, "do http request")
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

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

// Client is an abstraction of a solr client e.g. standard, json api or solr cloud (v2 api)
type Client interface {
	Query(ctx context.Context, collection string, query *Query) (*QueryResponse, error)
}

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

// Query calls the query endpoint with the provided query
func (c *JSONClient) Query(
	ctx context.Context,
	collection string,
	q *Query,
) (*QueryResponse, error) {
	urlStr := fmt.Sprintf("%s/solr/%s/query", c.baseURL, collection)
	theURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "build request url")
	}

	bq, err := q.BuildJSON()
	if err != nil {
		return nil, errors.Wrap(err, "build query body")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, theURL.String(), bytes.NewReader(bq))
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

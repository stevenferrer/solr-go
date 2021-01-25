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

// Client is an abstraction of a solr client
// e.g. standard, json api or solr cloud (v2 api)
type Client interface {
	Query(ctx context.Context, collection string, query *QueryBuilder) (*QueryResponse, error)
}

// JSONClient is a solr client that uses the JSON API
type JSONClient struct {
	proto      string
	host       string
	port       int
	httpClient *http.Client
}

var _ Client = (*JSONClient)(nil)

// NewJSONClient returns a new JSONClient
func NewJSONClient(host string, port int) *JSONClient {
	return &JSONClient{
		proto:      "http",
		host:       host,
		port:       port,
		httpClient: http.DefaultClient,
	}
}

// WithHTTPClient overrides the default http client
func (c *JSONClient) WithHTTPClient(httpClient *http.Client) *JSONClient {
	c.httpClient = httpClient
	return c
}

// Query is used for searching/querying
func (c *JSONClient) Query(ctx context.Context, collection string,
	qb *QueryBuilder) (*QueryResponse, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/query",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "build request url")
	}

	q, err := qb.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build query")
	}

	var b []byte
	b, err = json.Marshal(q)
	if err != nil {
		return nil, errors.Wrap(err, "marshal query")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, theURL.String(), bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do request")
	}

	var resp QueryResponse
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Client is a contract for querying with solr via JSON request API
type Client interface {
	Query(ctx context.Context, collection string,
		query map[string]interface{}) (*Response, error)
}

type client struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewClient is a factory for query client
func NewClient(host string, port int) Client {
	proto := "http"
	return &client{
		host:  host,
		port:  port,
		proto: proto,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// NewCustomClient is a factory for JSON query client with custom http client
func NewCustomClient(host string, port int, httpClient *http.Client) Client {
	proto := "http"
	return &client{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c client) Query(ctx context.Context, collection string,
	query map[string]interface{}) (*Response, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/query",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(query)
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

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

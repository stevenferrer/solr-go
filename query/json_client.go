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
	"github.com/stevenferrer/solr-go/types"
)

// JSONClient is a contract for interacting with Apache Solr JSON Request API
type JSONClient interface {
	Query(ctx context.Context, collection string, m types.M) (*Response, error)
}

type jsonClient struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewJSONClient is a factory for JSON query client
func NewJSONClient(host string, port int) JSONClient {
	proto := "http"
	return &jsonClient{
		host:  host,
		port:  port,
		proto: proto,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// NewJSONClientWithHTTPClient is a factory for JSON query client with custom http client
func NewJSONClientWithHTTPClient(host string, port int, httpClient *http.Client) JSONClient {
	proto := "http"
	return &jsonClient{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c jsonClient) Query(ctx context.Context, collection string, m types.M) (*Response, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/query",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "marshal m")
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

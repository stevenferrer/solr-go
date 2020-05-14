package helios

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gojek/heimdall/httpclient"
	"github.com/pkg/errors"
)

// Client is the contract for interacting with solr
type Client interface {
	// Select is used for querying. It accepts the request as JSON
	Select(ctx context.Context, collection string, request []byte) (*SelectResponse, error)
	// Update is used for indexing new data. It accepts the request as JSON
	Update(ctx context.Context, collection string, request []byte) (*UpdateResponse, error)
}

// client is the default implementation of solr client
type client struct {
	host string
	port int

	httpClient *httpclient.Client
}

// NewClient - a factory for creating default Client
func NewClient(host string, port int) Client {
	timeout := time.Second * 60
	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
	)

	return &client{
		host:       host,
		port:       port,
		httpClient: httpClient,
	}
}

func (c *client) Select(ctx context.Context, collection string, request []byte) (*SelectResponse, error) {
	theURL, err := url.Parse(
		fmt.Sprintf("http://%s:%d/solr/%s/query", c.host, c.port, collection),
	)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	headers := http.Header{
		"content-type": []string{"application/json"},
	}

	httpResp, err := c.httpClient.Post(
		theURL.String(),
		bytes.NewReader(request),
		headers,
	)
	if err != nil {
		return nil, errors.Wrap(err, "http post")
	}

	var response SelectResponse
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decode json response")
	}

	return &response, nil
}

func (c *client) Update(ctx context.Context, collection string, request []byte) (*UpdateResponse, error) {
	theURL, err := url.Parse(
		fmt.Sprintf("http://%s:%d/solr/%s/update?commit=true", c.host, c.port, collection),
	)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	headers := http.Header{
		"content-type": []string{"application/json"},
	}

	httpResp, err := c.httpClient.Post(
		theURL.String(),
		bytes.NewReader(request),
		headers,
	)
	if err != nil {
		return nil, errors.Wrap(err, "http post")
	}

	var response UpdateResponse
	err = json.NewDecoder(httpResp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decode json response")
	}

	return &response, nil
}

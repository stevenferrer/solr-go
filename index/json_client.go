package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// JSONClient is a contract for interacting with index API via JSON
type JSONClient interface {
	// AddDocs index multiple documents.
	// Note: make sure that "docs" is an array
	// TODO: improve this interface
	AddDocs(ctx context.Context, collection string, docs interface{}) error
	// UpdateCmds send multiple update commands
	UpdateCommands(ctx context.Context, collection string, commands ...Commander) error

	Commit(ctx context.Context, collection string) error
}

type jsonClient struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewJSONClient is a factory or JSON index client
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

// NewJSONClientWithHTTPClient is a factory or JSON index client with custom http client
func NewJSONClientWithHTTPClient(host string, port int, httpClient *http.Client) JSONClient {
	proto := "http"
	return &jsonClient{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c jsonClient) buildURL(collection string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	return u, nil
}

func (c jsonClient) AddDocs(ctx context.Context, collection string, docs interface{}) error {
	theURL, err := c.buildURL(collection)
	if err != nil {
		return errors.Wrap(err, "build url")
	}

	var b []byte
	b, err = json.Marshal(docs)
	if err != nil {
		return errors.Wrap(err, "marshal docs")
	}

	return c.update(ctx, theURL.String(), b)
}

func (c jsonClient) UpdateCommands(ctx context.Context, collection string, commands ...Commander) error {
	if len(commands) == 0 {
		return nil
	}

	theURL, err := c.buildURL(collection)
	if err != nil {
		return errors.Wrap(err, "build url")
	}

	cmdStrs := []string{}
	for _, cmd := range commands {
		cmdStr, err := cmd.Command()
		if err != nil {
			return errors.Wrap(err, "format commad")
		}

		cmdStrs = append(cmdStrs, cmdStr)
	}

	cmdBody := "{" + strings.Join(cmdStrs, ",") + "}"
	return c.update(ctx, theURL.String(), []byte(cmdBody))
}

func (c jsonClient) update(ctx context.Context, urlStr string, body []byte) error {
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urlStr, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}

	return c.doReq(ctx, httpReq)
}

func (c jsonClient) doReq(ctx context.Context, httpReq *http.Request) error {
	httpReq.Header.Add("content-type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "http do request")
	}

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil
}

func (c jsonClient) Commit(ctx context.Context, collection string) error {
	theURL, err := c.buildURL(collection)
	if err != nil {
		return errors.Wrap(err, "build url")
	}

	q := theURL.Query()
	q.Add("commit", "true")
	theURL.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodGet, theURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "new http request")
	}

	return c.doReq(ctx, httpReq)
}

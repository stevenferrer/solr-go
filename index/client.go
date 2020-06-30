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

// Client is a contract for interacting with index API via JSON API
type Client interface {
	// AddDocuments adds documents to the index
	AddDocuments(ctx context.Context, collection string, docs *Docs) error
	// SendCommands sends update commands
	SendCommands(ctx context.Context, collection string, commands ...Commander) error
	// Commit sends a commit request
	Commit(ctx context.Context, collection string) error
}

type client struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewClient is a factory or JSON index client
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

// NewCustomClient is a factory or JSON index client with custom options
func NewCustomClient(host string, port int, httpClient *http.Client) Client {
	proto := "http"
	return &client{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c client) buildURL(collection string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	return u, nil
}

func (c client) AddDocuments(ctx context.Context, collection string, docs *Docs) error {
	theURL, err := c.buildURL(collection)
	if err != nil {
		return errors.Wrap(err, "build url")
	}

	var b []byte
	b, err = docs.Marshal()
	if err != nil {
		return errors.Wrap(err, "marshal docs")
	}

	return c.update(ctx, theURL.String(), b)
}

func (c client) SendCommands(ctx context.Context, collection string, commands ...Commander) error {
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
			return errors.Wrap(err, "build commad")
		}

		cmdStrs = append(cmdStrs, cmdStr)
	}

	cmdBody := "{" + strings.Join(cmdStrs, ",") + "}"
	return c.update(ctx, theURL.String(), []byte(cmdBody))
}

func (c client) update(ctx context.Context, urlStr string, body []byte) error {
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urlStr, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}

	return c.doReq(ctx, httpReq)
}

func (c client) doReq(ctx context.Context, httpReq *http.Request) error {
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

func (c client) Commit(ctx context.Context, collection string) error {
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

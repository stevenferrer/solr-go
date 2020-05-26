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
	// AddSingle add a single document
	AddSingle(ctx context.Context, collection string, doc interface{}) error
	// AddMultiple add multiple documents
	AddMultiple(ctx context.Context, collection string, docs interface{}) error
	// UpdateCmds send multiple update commands
	UpdateCommands(ctx context.Context, collection string, commands ...Commander) error
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

func (c jsonClient) AddSingle(ctx context.Context, collection string, doc interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update/json/docs?commit=true",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(doc)
	if err != nil {
		return errors.Wrap(err, "marshal doc")
	}

	return c.doUpdt(ctx, theURL.String(), b)
}

func (c jsonClient) AddMultiple(ctx context.Context, collection string, docs interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update?commit=true",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(docs)
	if err != nil {
		return errors.Wrap(err, "marshal docs")
	}

	return c.doUpdt(ctx, theURL.String(), b)
}

func (c jsonClient) UpdateCommands(ctx context.Context, collection string, commands ...Commander) error {
	if len(commands) == 0 {
		return nil
	}

	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update?commit=true",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return errors.Wrap(err, "parse url")
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
	return c.doUpdt(ctx, theURL.String(), []byte(cmdBody))
}

func (c jsonClient) doUpdt(ctx context.Context, urlStr string, body []byte) error {
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urlStr, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
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

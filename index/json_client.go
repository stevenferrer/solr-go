package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// JSONClient is a JSON index client
type JSONClient interface {
	// AddSingle add a single document
	AddSingle(ctx context.Context, coll string, doc interface{}) error
	// AddMultiple add multiple documents
	AddMultiple(ctx context.Context, coll string, docs interface{}) error
	// UpdateCmds send multiple update commands
	UpdateCmds(ctx context.Context, coll string, cmds ...Cmd) error
}

type jsonClient struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewJSONClient is a factory or JSON index client
func NewJSONClient(host string, port int, httpClient *http.Client) JSONClient {
	proto := "http"
	return &jsonClient{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c jsonClient) AddSingle(ctx context.Context, coll string, doc interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update/json/docs?commit=true",
		c.proto, c.host, c.port, coll))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var body []byte
	body, err = json.Marshal(doc)
	if err != nil {
		return errors.Wrap(err, "marshal doc")
	}

	return c.doUpdt(ctx, theURL.String(), body)
}

func (c jsonClient) AddMultiple(ctx context.Context, coll string, docs interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update?commit=true",
		c.proto, c.host, c.port, coll))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var body []byte
	body, err = json.Marshal(docs)
	if err != nil {
		return errors.Wrap(err, "marshal docs")
	}

	return c.doUpdt(ctx, theURL.String(), body)
}

func (c jsonClient) UpdateCmds(ctx context.Context, coll string, cmds ...Cmd) error {
	if len(cmds) == 0 {
		return nil
	}

	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/update?commit=true",
		c.proto, c.host, c.port, coll))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	cmdStrs := []string{}
	for _, cmd := range cmds {
		cmdStr, err := cmd.ToCmd()
		if err != nil {
			return errors.Wrap(err, "format commad")
		}

		cmdStrs = append(cmdStrs, cmdStr)
	}

	cmdBody := "{" + strings.Join(cmdStrs, ",") + "}"
	return c.doUpdt(ctx, theURL.String(), []byte(cmdBody))
}

func (c jsonClient) doUpdt(ctx context.Context, theURL string, body []byte) error {
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, theURL, bytes.NewReader(body))
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
		// return resp.Error
		return resp.Error
	}

	return nil
}

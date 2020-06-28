package config

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

// Client is a config API client
type Client interface {
	GetConfig(ctx context.Context, collection string) (*Response, error)
	SendCommands(ctx context.Context, collection string, commands ...Commander) error
}

type client struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// New is a factory for config client
func New(host string, port int) Client {
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

// NewWithHTTPClient is a factory for config client
func NewWithHTTPClient(host string, port int, httpClient *http.Client) Client {
	proto := "http"
	return &client{
		host:       host,
		port:       port,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c *client) GetConfig(ctx context.Context, collection string) (*Response, error) {
	theURL, err := c.buildURL(collection)
	if err != nil {
		return nil, errors.Wrap(err, "build url")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, theURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}

	return c.do(httpReq)
}

func (c *client) SendCommands(ctx context.Context, collection string, commands ...Commander) error {
	if len(commands) == 0 {
		return nil
	}

	theURL, err := c.buildURL(collection)
	if err != nil {
		return errors.Wrap(err, "build url")
	}

	// build commands
	commandStrs := []string{}
	for _, command := range commands {
		commandStr, err := command.Command()
		if err != nil {
			return errors.Wrap(err, "build commad")
		}

		commandStrs = append(commandStrs, commandStr)
	}

	// send commands to solr
	requestBody := "{" + strings.Join(commandStrs, ",") + "}"

	return c.sendCommands(ctx, theURL.String(), []byte(requestBody))
}

func (c *client) buildURL(collection string) (*url.URL, error) {
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/config",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	return u, nil
}

func (c *client) sendCommands(ctx context.Context, urlStr string, body []byte) error {
	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost, urlStr, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}

	_, err = c.do(httpReq)
	if err != nil {
		return errors.Wrap(err, "send commands")
	}

	return err
}

func (c *client) do(httpReq *http.Request) (*Response, error) {
	httpReq.Header.Add("content-type", "application/json")
	httpResp, err := c.httpClient.Do(httpReq)
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

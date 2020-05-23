package suggester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Client is the suggester client
type Client interface {
	Suggest(context.Context, Request) (*Response, error)
}

type client struct {
	host       string
	port       int
	proto      string
	endpoint   string
	httpClient *http.Client
}

// NewClient is a factory for suggester Client
func NewClient(host string, port int,
	endpoint string, httpClient *http.Client) Client {
	proto := "http"

	return &client{
		host:       host,
		port:       port,
		endpoint:   endpoint,
		proto:      proto,
		httpClient: httpClient,
	}
}

func (c client) Suggest(ctx context.Context, req Request) (*Response, error) {
	if req.Params.Query == "" {
		return nil, errors.New("query is required")
	}

	params := buildParams(req.Params)

	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/%s?%s", c.proto,
		c.host, c.port, req.Collection, c.endpoint, strings.Join(params, "&")))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	httpReq, err := http.NewRequestWithContext(ctx,
		http.MethodGet, theURL.String(), nil)
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

func buildParams(p Params) []string {
	params := []string{}

	params = append(params, "suggest=true",
		fmt.Sprintf("suggest.q=%s", p.Query))

	for _, dict := range p.Dictionaries {
		params = append(params, fmt.Sprintf("suggest.dictionary=%s", dict))
	}

	if p.Count > 0 {
		params = append(params, fmt.Sprintf("suggest.count=%d", p.Count))
	}

	if p.Cfq != "" {
		params = append(params, fmt.Sprintf("suggest.cfg=%s", p.Cfq))
	}

	if p.Build {
		params = append(params, "suggest.build=true")
	}

	if p.Reload {
		params = append(params, "suggest.reload=true")
	}

	if p.BuildAll {
		params = append(params, "suggest.buildAll=true")
	}

	if p.ReloadAll {
		params = append(params, "suggest.reloadAll=true")
	}

	return params
}

package suggester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Client is the suggester client
type Client interface {
	Suggest(ctx context.Context, collection string, params Params) (*Response, error)
}

// Params is the suggester parameters
type Params struct {
	// Dictionaries is the name of the dictionary
	// component configured in the search component
	Dictionaries []string

	// Query is the query to use for suggestion lookups
	Query string

	// Count is the number of suggestions to return
	Count int

	// ContextFilterQuery  is a context filter query used to filter
	// suggestsion based on the context field, if supported by the suggester
	Cfq string

	// Build If true, it will build the suggester index
	Build,

	// Reload If true, it will reload the suggester index.
	Reload,

	// BuildAll If true, it will build all suggester indexes.
	BuildAll,

	// ReloadAll If true, it will reload all suggester indexes.
	ReloadAll bool
}

type client struct {
	host       string
	port       int
	proto      string
	endpoint   string
	httpClient *http.Client
}

// NewClient is a factory for suggester
// Client with default configurations
func NewClient(host string, port int) Client {
	proto := "http"
	return &client{
		host:     host,
		port:     port,
		endpoint: "suggest",
		proto:    proto,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// NewCustomClient is a factory for suggester Client with custom configurations
func NewCustomClient(host string, port int,
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

func (c client) Suggest(ctx context.Context, collection string, params Params) (*Response, error) {
	if params.Query == "" {
		return nil, errors.New("query is required")
	}

	paramList := buildParams(params)

	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/%s?%s", c.proto,
		c.host, c.port, collection, c.endpoint, strings.Join(paramList, "&")))
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

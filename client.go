package solr

import (
	"net/http"

	"github.com/sf9v/solr-go/config"
	"github.com/sf9v/solr-go/index"
	"github.com/sf9v/solr-go/query"
	"github.com/sf9v/solr-go/schema"
	"github.com/sf9v/solr-go/suggester"
)

// Client is a unified solr client
type Client interface {
	Index() index.Client
	Query() query.Client
	Schema() schema.Client
	Suggester() suggester.Client
	Config() config.Client
}

type client struct {
	indexClient     index.Client
	queryClient     query.Client
	schemaClient    schema.Client
	suggesterClient suggester.Client
	configClient    config.Client
}

// NewClient is a factory for solr unified client
func NewClient(host string, port int) Client {
	return &client{
		indexClient:     index.NewClient(host, port),
		queryClient:     query.NewClient(host, port),
		schemaClient:    schema.NewClient(host, port),
		suggesterClient: suggester.NewClient(host, port),
		configClient:    config.NewClient(host, port),
	}
}

// NewCustomClient is a factory for solr unified client with custom options
func NewCustomClient(host string, port int, httpClient *http.Client) Client {
	return &client{
		indexClient:     index.NewCustomClient(host, port, httpClient),
		queryClient:     query.NewCustomClient(host, port, httpClient),
		schemaClient:    schema.NewCustomClient(host, port, httpClient),
		suggesterClient: suggester.NewCustomClient(host, port, "suggest", httpClient),
		configClient:    config.NewCustomClient(host, port, httpClient),
	}
}

func (c *client) Index() index.Client {
	return c.indexClient
}

func (c *client) Query() query.Client {
	return c.queryClient
}

func (c *client) Schema() schema.Client {
	return c.schemaClient
}

func (c *client) Suggester() suggester.Client {
	return c.suggesterClient
}

func (c *client) Config() config.Client {
	return c.configClient
}

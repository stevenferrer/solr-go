package solr

import (
	"github.com/stevenferrer/solr-go/config"
	"github.com/stevenferrer/solr-go/index"
	"github.com/stevenferrer/solr-go/query"
	"github.com/stevenferrer/solr-go/schema"
	"github.com/stevenferrer/solr-go/suggester"
)

// Client is a unified solr client
type Client interface {
	Index() index.JSONClient
	Query() query.JSONClient
	Schema() schema.Client
	Suggester() suggester.Client
	Config() config.Client
}

type client struct {
	indexClient     index.JSONClient
	queryClient     query.JSONClient
	schemaClient    schema.Client
	suggesterClient suggester.Client
	configClient    config.Client
}

// NewClient is a factory for solr Client
func NewClient(host string, port int) Client {
	return &client{
		indexClient:     index.NewJSONClient(host, port),
		queryClient:     query.NewJSONClient(host, port),
		schemaClient:    schema.NewClient(host, port),
		suggesterClient: suggester.NewClient(host, port),
		configClient:    config.New(host, port),
	}
}

func (c *client) Index() index.JSONClient {
	return c.indexClient
}

func (c *client) Query() query.JSONClient {
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

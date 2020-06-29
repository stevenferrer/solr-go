package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	solr "github.com/stevenferrer/solr-go"
)

func TestNewClient(t *testing.T) {
	// only for covering
	host := "localhost"
	port := 8983

	// only for covering
	client := solr.NewClient(host, port)

	assert.NotNil(t, client.Index())
	assert.NotNil(t, client.Query())
	assert.NotNil(t, client.Schema())
	assert.NotNil(t, client.Suggester())
	assert.NotNil(t, client.Config())
}

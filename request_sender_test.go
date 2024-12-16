package solr_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/public-safety-canada/solr-go"
)

func TestRequestSender(t *testing.T) {
	rs := solr.NewDefaultRequestSender().
		WithHTTPClient(http.DefaultClient).
		WithBasicAuth("solr", "SolrRocks")

	ctx := context.Background()
	_, err := rs.SendRequest(ctx, "", "", "", nil)
	assert.Error(t, err)

	_, err = rs.SendRequest(ctx, ":", "", "", nil)
	assert.Error(t, err)
}

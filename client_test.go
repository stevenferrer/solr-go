package solr_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"

	"github.com/sf9v/solr-go"
)

func TestJSONClient(t *testing.T) {
	defer gock.Off()

	a := assert.New(t)
	baseURL := "https://solr.example.com"
	collection := "products"

	gock.New(baseURL).
		Post("/solr/" + collection + "/query").
		MatchType("json").
		BodyString(`{"query":"{!dismax}apple pie"}`).
		Reply(http.StatusOK).
		JSON(solr.QueryResponse{})

	client := solr.NewJSONClient(baseURL).
		WithHTTPClient(&http.Client{
			Timeout: 60 * time.Second,
		})
	q := solr.NewQuery().
		WithQueryParser(
			solr.NewDisMaxQueryParser("apple pie"),
		)

	_, err := client.Query(context.Background(), collection, q)
	a.NoError(err)

	a.True(gock.IsDone())
}

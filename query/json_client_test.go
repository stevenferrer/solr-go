package query_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/stevenferrer/solr-go"
	"github.com/stevenferrer/solr-go/query"
)

func TestJSONClient(t *testing.T) {
	ctx := context.Background()
	collection := "techproducts"
	host := "localhost"
	port := 8983
	timeout := time.Second * 60

	// Test examples are extracted from Apache Solr website
	// https://lucene.apache.org/solr/guide/8_5/json-query-dsl.html

	t.Run("errors", func(t *testing.T) {
		t.Run("query error", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/query-error")
			require.NoError(t, err)
			defer rec.Stop()

			client := query.NewJSONClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			_, err = client.Query(ctx, collection, M{
				"curry": "",
			})
			a.Error(err)
		})

		t.Run("parse url error", func(t *testing.T) {
			a := assert.New(t)
			client := query.NewJSONClient("wtf:\\//wtf::", port, &http.Client{})
			_, err := client.Query(ctx, collection, M{})
			a.Error(err)
		})

		t.Run("request error", func(t *testing.T) {
			a := assert.New(t)
			client := query.NewJSONClient(host, 1234, &http.Client{})
			_, err := client.Query(ctx, collection, M{})
			a.Error(err)
		})
	})

	t.Run("simple query string", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/simple-query-string")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": "name:iPod",
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("local-params string", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/local-params-string")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": "{!lucene df=name v=iPod}",
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("full expanded JSON object", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/fully-expanded-json-object")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": M{
				"lucene": M{
					"df":    "name",
					"query": "iPod",
				},
			},
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("nested boost queries", func(t *testing.T) {
		t.Run("mix of fully-expanded and local-params", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/boost-mix-fully-expanded-and-local-params")
			require.NoError(t, err)
			defer rec.Stop()

			client := query.NewJSONClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			var resp *query.Response
			resp, err = client.Query(ctx, collection, M{
				"query": M{
					"boost": M{
						"query": "{!lucene df=name}iPod",
						"b":     "log(popularity)",
					},
				},
			})
			a.NoError(err)
			a.NotNil(resp)
		})

		t.Run("all queries fully-expanded as JSON", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/boost-all-fully-expanded-as-json")
			require.NoError(t, err)
			defer rec.Stop()

			client := query.NewJSONClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			var resp *query.Response
			resp, err = client.Query(ctx, collection, M{
				"query": M{
					"boost": M{
						"query": M{
							"lucene": M{
								"df":    "name",
								"query": "iPod",
							},
						},
						"b": "log(popularity)",
					},
				},
			})
			a.NoError(err)
			a.NotNil(resp)
		})
	})

	t.Run("nested boolean query", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/boolean-query")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": M{
				"bool": M{
					"must": []M{
						{"lucene": M{"df": "name", "query": "iPod"}},
					},
					"must_not": []M{
						{"frange": M{"l": "0", "u": "5", "query": "popularity"}},
					},
				},
			},
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("filter query", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/filter-query")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": M{
				"bool": M{
					"must_not": "{!frange l=0 u=5}popularity",
				},
			},
			"filter": []string{"name:iPod"},
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("additional query", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/additional-query")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"queries": M{
				"electronic": M{"field": M{"f": "cat", "query": "electronics"}},
				"manufacturers": []interface{}{
					"manu:apple",
					M{"field": M{"f": "manu", "query": "belkin"}},
				},
			},
			"query": "{!v=$electronic}",
		})
		a.NoError(err)
		a.NotNil(resp)
	})

	t.Run("tagging in json query dsl", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/tagging-in-json-query-dsl")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, collection, M{
			"query": "*:*",
			"filter": []M{
				{"#titleTag": "name:Solr"},
				{"#inStockTag": "inStock:true"},
			},
		})
		a.NoError(err)
		a.NotNil(resp)

	})
}

package query_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/solr-go/query"
	. "github.com/stevenferrer/solr-go/types"
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

			client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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
			client := query.NewJSONClientWithHTTPClient("wtf:\\//wtf::", port, &http.Client{})
			_, err := client.Query(ctx, collection, M{})
			a.Error(err)
		})

		t.Run("request error", func(t *testing.T) {
			a := assert.New(t)
			client := query.NewJSONClientWithHTTPClient(host, 1234, &http.Client{})
			_, err := client.Query(ctx, collection, M{})
			a.Error(err)
		})
	})

	t.Run("simple query string", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/simple-query-string")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

			client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

			client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
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

	t.Run("multi-select facet", func(t *testing.T) {
		a := assert.New(t)

		rec, err := recorder.New("fixtures/multi-select-facet")
		require.NoError(t, err)
		defer rec.Stop()

		client := query.NewJSONClientWithHTTPClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		var resp *query.Response
		resp, err = client.Query(ctx, "multi-select", M{
			"queries": M{
				"child.query": "scope:sku",
				"child.fq": []string{
					"{!tag=color}color:black OR color:blue",
					"{!tag=size}size:L OR size:M",
				},
				"parent.fq": "scope:product",
			},
			"query": "{!parent tag=top filters=$child.fq which=$parent.fq v=$child.query}",
			"filter": []string{
				"{!tag=top}category:clothes",
			},
			"fields": "*",
			"facet": M{
				"colors": M{
					"domain": M{
						"excludeTags": "top",
						"filter": []string{
							"{!filters param=$child.fq  excludeTags=color v=$child.query}",
							"{!child of=$parent.fq filters=$fq v=$parent.fq}",
						},
					},
					"type":  "terms",
					"field": "color",
					"limit": -1,
					"facet": M{
						"productCount": "uniqueBlock(_root_)",
					},
				},

				"sizes": M{
					"domain": M{
						"excludeTags": "top",
						"filter": []string{
							"{!filters param=$child.fq excludeTags=size v=$child.query}",
							"{!child of=$parent.fq filters=$fq v=$parent.fq}",
						},
					},
					"type":  "terms",
					"field": "size",
					"limit": -1,
					"facet": M{
						"productCount": "uniqueBlock(_root_)",
					},
				},

				"categories": M{
					"type":  "terms",
					"field": "category",
					"limit": -1,
					"facet": M{
						"productCount": "uniqueBlock(_root_)",
					},
				},

				"brands": M{
					"type":  "terms",
					"field": "brand",
					"limit": -1,
					"facet": M{
						"productCount": "uniqueBlock(_root_)",
					},
				},
				"product_types": M{
					"type":  "terms",
					"field": "product_type",
					"limit": -1,
					"facet": M{
						"productCount": "uniqueBlock(_root_)",
					},
				},
			},
		})
		a.NoError(err)
		a.NotNil(resp)

		a.Len(resp.Facets, 6)
	})
}

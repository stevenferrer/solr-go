package solr_test

import (
	"testing"

	"github.com/sf9v/solr-go"
	"github.com/stretchr/testify/assert"
)

func TestFacets(t *testing.T) {
	t.Run("terms facet", func(t *testing.T) {
		nestedFacet := solr.NewQueryFacet("high_popularity").
			WithQuery("popularity:[8 TO 10]")
		facet := solr.NewTermsFacet("categories").
			WithField("cat").WithLimit(5).
			AddToFacet("average_price", "avg(price)").
			AddNestedFacet(nestedFacet)
		got := facet.BuildFacet()

		expect := solr.M{
			"type":  "terms",
			"field": "cat",
			"limit": 5,
			"facet": solr.M{
				"average_price": "avg(price)",
				"high_popularity": solr.M{
					"type": "query",
					"q":    "popularity:[8 TO 10]",
				},
			},
		}

		assert.Equal(t, "categories", facet.Name())
		assert.Equal(t, expect, got)
	})

	t.Run("query facet", func(f *testing.T) {
		nestedFacet := solr.NewTermsFacet("categories").
			WithField("cat").WithLimit(5)
		facet := solr.NewQueryFacet("high_popularity").
			WithQuery("popularity:[8 TO 10]").
			AddToFacet("average_price", "avg(price)").
			AddNestedFacet(nestedFacet)
		got := facet.BuildFacet()

		expect := solr.M{
			"type": "query",
			"q":    "popularity:[8 TO 10]",
			"facet": solr.M{
				"average_price": "avg(price)",
				"categories": solr.M{
					"type":  "terms",
					"field": "cat",
					"limit": 5,
				},
			},
		}

		assert.Equal(t, "high_popularity", facet.Name())
		assert.Equal(t, expect, got)
	})
}

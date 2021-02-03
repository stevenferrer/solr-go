package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestFacets(t *testing.T) {
	t.Run("terms facet", func(t *testing.T) {
		queryFacet := solr.NewQueryFacet("high_popularity").
			Query("popularity:[8 TO 10]")
		facet := solr.NewTermsFacet("categories").
			Field("cat").Limit(5).Offset(1).Sort("price asc").
			AddFacets(queryFacet).
			AddToFacet("average_price", "avg(price)").
			AddToDomain("excludeTags", "top")
		got := facet.BuildFacet()

		expect := solr.M{
			"type":   "terms",
			"field":  "cat",
			"limit":  5,
			"offset": 1,
			"sort":   "price asc",
			"facet": solr.M{
				"average_price": "avg(price)",
				"high_popularity": solr.M{
					"type": "query",
					"q":    "popularity:[8 TO 10]",
				},
			},
			"domain": solr.M{
				"excludeTags": "top",
			},
		}

		assert.Equal(t, "categories", facet.Name())
		assert.Equal(t, expect, got)
	})

	t.Run("query facet", func(f *testing.T) {
		termsFacet := solr.NewTermsFacet("categories").
			Field("cat").Limit(5)
		facet := solr.NewQueryFacet("high_popularity").
			Query("popularity:[8 TO 10]").
			AddFacet(termsFacet).
			AddToFacet("average_price", "avg(price)")
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

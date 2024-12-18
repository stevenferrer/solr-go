package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/public-safety-canada/solr-go"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)
	got := solr.NewQuery(solr.NewDisMaxQueryParser().
		Query("'solr rocks'").BuildParser()).
		Queries(solr.M{
			"query_filters": []solr.M{
				{
					"#size_tag": solr.M{
						"field": solr.M{
							"f":     "size",
							"query": "XL",
						},
					},
				},
				{
					"#color_tag": solr.M{
						"field": solr.M{
							"f":     "color",
							"query": "Red",
						},
					},
				},
			},
		}).
		Facets(
			solr.NewTermsFacet("categories").
				Field("cat").Limit(10),
			solr.NewQueryFacet("high_popularity").
				Query("popularity:[8 TO 10]"),
		).
		Sort("score").
		Offset(1).
		Limit(10).
		Filters("inStock:true").
		Fields("name", "price").
		BuildQuery()

	expect := solr.M{
		"facet": solr.M{
			"categories":      solr.M{"field": "cat", "limit": 10, "type": "terms"},
			"high_popularity": solr.M{"q": "popularity:[8 TO 10]", "type": "query"},
		},
		"fields": []string{"name", "price"},
		"filter": []string{"inStock:true"},
		"limit":  10,
		"offset": 1,
		"queries": solr.M{
			"query_filters": []solr.M{
				{"#size_tag": solr.M{"field": solr.M{"f": "size", "query": "XL"}}},
				{"#color_tag": solr.M{"field": solr.M{"f": "color", "query": "Red"}}},
			},
		},
		"query": "{!dismax v='solr rocks'}",
		"sort":  "score",
	}

	a.Equal(expect, got)



	a = assert.New(t)
	got = solr.NewQuery(solr.NewExtendedDisMaxQueryParser().
		Query("'solr rocks'").BuildParser()).
		Queries(solr.M{
			"query_filters": []solr.M{
				{
					"#size_tag": solr.M{
						"field": solr.M{
							"f":     "size",
							"query": "XL",
						},
					},
				},
				{
					"#color_tag": solr.M{
						"field": solr.M{
							"f":     "color",
							"query": "Red",
						},
					},
				},
			},
		}).
		Facets(
			solr.NewTermsFacet("categories").
				Field("cat").Limit(10),
			solr.NewQueryFacet("high_popularity").
				Query("popularity:[8 TO 10]"),
		).
		Sort("score").
		Offset(1).
		Limit(10).
		Filters("inStock:true").
		Fields("name", "price").
		BuildQuery()

	expect = solr.M{
		"facet": solr.M{
			"categories":      solr.M{"field": "cat", "limit": 10, "type": "terms"},
			"high_popularity": solr.M{"q": "popularity:[8 TO 10]", "type": "query"},
		},
		"fields": []string{"name", "price"},
		"filter": []string{"inStock:true"},
		"limit":  10,
		"offset": 1,
		"queries": solr.M{
			"query_filters": []solr.M{
				{"#size_tag": solr.M{"field": solr.M{"f": "size", "query": "XL"}}},
				{"#color_tag": solr.M{"field": solr.M{"f": "color", "query": "Red"}}},
			},
		},
		"query": "{!edismax v='solr rocks'}",
		"sort":  "score",
	}

	a.Equal(expect, got)
}

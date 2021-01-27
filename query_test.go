package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)
	got, err := solr.NewQuery().
		WithQueryParser(
			solr.NewDisMaxQueryParser("solr rocks"),
		).
		WithQueries(solr.M{
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
		WithFacets(
			solr.NewTermsFacet("categories").
				WithField("cat").WithLimit(10),
		).
		WithSort("score").
		WithOffset(1).
		WithLimit(10).
		WithFilter("inStock:true").
		WithFields("name price").
		BuildJSON()
	a.NoError(err)
	expect := `{"facet":{"categories":{"field":"cat","limit":10,"type":"terms"}},"fields":"name price","filter":"inStock:true","limit":10,"offset":1,"queries":{"query_filters":[{"#size_tag":{"field":{"f":"size","query":"XL"}}},{"#color_tag":{"field":{"f":"color","query":"Red"}}}]},"query":"{!dismax}solr rocks","sort":"score"}`
	a.Equal(expect, string(got))
}

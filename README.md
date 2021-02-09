![Github Actions](https://github.com/sf9v/solr-go/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/sf9v/solr-go/badge.svg?branch=main)](https://coveralls.io/github/sf9v/solr-go?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/sf9v/solr-go)](https://goreportcard.com/report/github.com/sf9v/solr-go)

# Solr-Go

A [Solr](https://lucene.apache.org/solr) client for [Go](https://golang.org/).

## Supported APIs

- [Collections API](https://lucene.apache.org/solr/guide/8_8/collections-api.html)
- [Query API](https://lucene.apache.org/solr/guide/8_8/json-request-api.html)
  - [Facet API](https://lucene.apache.org/solr/guide/8_8/json-facet-api.html)
- [Update API](https://lucene.apache.org/solr/guide/8_8/uploading-data-with-index-handlers.html#uploading-data-with-index-handlers)
- [Schema API](https://lucene.apache.org/solr/guide/8_8/schema-api.html)
- [Config API](https://lucene.apache.org/solr/guide/8_8/config-api.html)
- [Suggester API](https://lucene.apache.org/solr/guide/8_8/suggester.html)

## Example

See [integration test](integration_test.go) for more examples.

```go
// Create a client
baseURL := "http://solr.example.com"
client := solr.NewJSONClient(baseURL)

// Create a query
query := solr.NewQuery().
    QueryParser(
        solr.NewDisMaxQueryParser().
            Query("'solr rocks'"),
    ).
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
    Fields("name", "price")

// Send the query
queryResponse, err := client.Query(context.Background(), "techproducts", query)
```

## Contributing

Any contributions are welcome!

## License

MIT

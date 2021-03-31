# Solr-Go

![Github Actions](https://github.com/sf9v/solr-go/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/sf9v/solr-go/badge.svg?branch=main)](https://coveralls.io/github/sf9v/solr-go?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/sf9v/solr-go)](https://goreportcard.com/report/github.com/sf9v/solr-go)
[![GoDoc](https://pkg.go.dev/badge/github.com/sf9v/solr-go)](https://pkg.go.dev/github.com/sf9v/solr-go)

A [Solr](https://lucene.apache.org/solr) client for [Go](https://golang.org/).

## Installation

```console
$ go get github.com/sf9v/solr-go
```

## Example

See [integration test](integration_test.go) for more examples.

```go
// Create a client
baseURL := "http://solr.example.com"
client := solr.NewJSONClient(baseURL)

// Create a query
query := solr.NewQuery(solr.NewDisMaxQueryParser().
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
    Fields("name", "price")

// Send the query
queryResponse, err := client.Query(context.Background(), "techproducts", query)
```

## Supported APIs

- [Collections API](https://solr.apache.org/guide/8_8/collections-api.html)
- [Core Admin API]()
- [Query API](https://solr.apache.org/guide/8_8/json-request-api.html)
  - [Facet API](https://solr.apache.org/guide/8_8/json-facet-api.html)
- [Update API](https://solr.apache.org/guide/8_8/uploading-data-with-index-handlers.html#uploading-data-with-index-handlers)
- [Schema API](https://solr.apache.org/guide/8_8/schema-api.html)
- [Config API](https://solr.apache.org/guide/8_8/config-api.html)
- [Suggester API](https://solr.apache.org/guide/8_8/suggester.html)

## Other supported features

- [Basic Auth](https://solr.apache.org/guide/8_8/basic-authentication-plugin.html#basic-authentication-plugin)

## Projects using it

- [Multi-select facet using Solr, Vue and Go](https://github.com/sf9v/multi-select-facet)

## Note

The current API is still evolving and will likely change before it hits version 1. Please use this library with caution.

## Contributing

Please feel free to improve this project by [opening an issue](https://github.com/sf9v/solr-go/issues/new) or by [making a pull-request](https://github.com/sf9v/solr-go/pulls).

## License

MIT

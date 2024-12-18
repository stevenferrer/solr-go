# Solr-Go

![Github Actions](https://github.com/public-safety-canada/solr-go/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/public-safety-canada/solr-go/badge.svg?branch=main)](https://coveralls.io/github/public-safety-canada/solr-go?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/public-safety-canada/solr-go)](https://goreportcard.com/report/github.com/public-safety-canada/solr-go)
[![GoDoc](https://pkg.go.dev/badge/github.com/public-safety-canada/solr-go)](https://pkg.go.dev/github.com/public-safety-canada/solr-go)

A [Solr](https://solr.apache.org/) client for [Go](https://golang.org/).

## Installation

```console
$ go get github.com/public-safety-canada/solr-go
```

## Example

```go
import (
    "github.com/public-safety-canada/solr-go"
)

...

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

See [integration test](integration_test.go) for a more complete example.

## Supported APIs

- [Collections API](https://solr.apache.org/guide/collections-api.html) - Create and delete collection.
- [Core Admin API](https://solr.apache.org/guide/coreadmin-api.html) - [Create](https://issues.apache.org/jira/browse/SOLR-7316), delete and check core status.
- [Query API](https://solr.apache.org/guide/json-request-api.html) - Query via the JSON request API.
  - [Facet API](https://solr.apache.org/guide/json-facet-api.html) - Terms and query facet.
- [Update API](https://solr.apache.org/guide/uploading-data-with-index-handlers.html#uploading-data-with-index-handlers) - JSON formatted index updates.
- [Schema API](https://solr.apache.org/guide/schema-api.html) - Modify schema fields, dynamic fields, copy fields and field types.
- [Config API](https://solr.apache.org/guide/config-api.html) - Modify config properties and update components.
- [Suggester API](https://solr.apache.org/guide/suggester.html) - Auto-suggest/type-ahead via suggester component.

## Supported Query Parsers

- [Standard Query Parser](https://solr.apache.org/guide/solr/latest/query-guide/standard-query-parser.html)
- [DisMax Query Parser](https://solr.apache.org/guide/solr/latest/query-guide/dismax-query-parser.html)
- [Extended DisMax (eDisMax) Query Parser](https://solr.apache.org/guide/solr/latest/query-guide/edismax-query-parser.html)

## Other features

- [Basic auth support](https://solr.apache.org/guide/basic-authentication-plugin.html#basic-authentication-plugin) - Interacting with a Solr server that uses the basic authentication plugin.

## Currently tested versions

- [Solr 9.7](https://solr.apache.org/guide/9_7/index.html)
- [Solr 8.11](https://solr.apache.org/guide/8_11/index.html)

## A (minor) word of caution

Keep in mind that this library is still evolving and will likely have some breaking changes until v1.0. We will try our best to keep the breaking changes minimal.

## Supporting the project

You can support the project in the following ways:

- Give it a [star](https://github.com/public-safety-canada/solr-go/stargazers), it's free!
- Write some tutorials
- Use it your projects

## Contributing

Please feel free to improve this project by [opening an issue](https://github.com/public-safety-canada/solr-go/issues/new) or by [making a pull-request](https://github.com/public-safety-canada/solr-go/pulls).

## License

MIT

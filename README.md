[![Go Report Card](https://goreportcard.com/badge/github.com/sf9v/solr-go)](https://goreportcard.com/report/github.com/sf9v/solr-go)

# Solr-Go

[Solr](https://lucene.apache.org/solr/) client for [Go](http://go.dev/). 

```go

import (
    "context"
    // Import the package that you need
    solrquery "github.com/sf9v/solr-go/query"
)

func main() {
    // Initialize the query client
    queryClient := solrquery.NewClient("localhost", 8983)

    // Start querying!
    queryResp, err := queryClient.Query(
        context.Background(),
        "techproducts", // name of your collection
        map[string]string{
            "query": "{!lucene df=name v=iPod}",
        },
    )
    ...
}
```

## Contents

- [Solr-Go](#solr-go)
  - [Contents](#contents)
  - [Notes](#notes)
  - [Projects using it](#projects-using-it)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Features](#features)
  - [Contributing](#contributing)

## Notes

* This is a *WORK IN-PROGRESS*, API might change a lot before *v1*
* I'm currently using this module in my projects
* Tested using [Solr 8.5](https://lucene.apache.org/solr/guide/8_5/)

## Projects using it

- [Multi-select facet using Solr, Vue and Go](https://github.com/sf9v/multi-select-facet)

## Installation

```console
$ go get github.com/sf9v/solr-go
```

## Usage

A detailed documentation shall follow after *v1*. For now you can start looking at the *tests* or *examples* inside each package directory.

* [Index API example](./index/examples/main.go)
* [Query API example](./query/example/main.go)
* [Schema API example](./schema/example/main.go)
* Suggester API example - TODO

## Features

- [x] [Schema API client](https://lucene.apache.org/solr/guide/8_5/schema-api.html) client
  - [x] [Modify schema](https://lucene.apache.org/solr/guide/8_5/schema-api.html#modify-the-schema)
  - [x] [Retrieve schema information](https://lucene.apache.org/solr/guide/8_5/schema-api.html#retrieve-schema-information)
  - [x] Example
- [ ] Index API
  - [x] [JSON client](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html)
    - [x] [Add multiple documents](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-multiple-json-documents)
    - [x] [Send update commands](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands) 
	- [x] Example
  - [ ] XML client ??
  - [ ] CSV client ??
- [x] [JSON query API client](https://lucene.apache.org/solr/guide/8_5/json-query-dsl.html)
  - [x] Facet
  - [x] Example
- [ ] [Standard query API client](https://lucene.apache.org/solr/guide/8_5/the-standard-query-parser.html#the-standard-query-parser)??
  - [ ] Example
- [x] [Suggester client](https://lucene.apache.org/solr/guide/8_5/suggester.html)
  - [ ] Example
- [x] [Config API client](https://lucene.apache.org/solr/guide/8_5/config-api.html)
  - [ ] Example
- [ ] [Collections API client](https://lucene.apache.org/solr/guide/8_5/collections-api.html)
  - [ ] Example
- [ ] [Configset API client](https://lucene.apache.org/solr/guide/8_5/configsets-api.html)
  - [ ] Example
- [x] Unified solr client
- [ ] SolrCloud support (V2 API)
- [ ] Basic auth support
- [ ] Documentation

## Contributing

This is a **work in-progress**, any contributions are very welcome!

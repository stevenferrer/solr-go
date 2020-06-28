[![Go Report Card](https://goreportcard.com/badge/github.com/stevenferrer/solr-go)](https://goreportcard.com/report/github.com/stevenferrer/solr-go)
[![CircleCI](https://circleci.com/gh/stevenferrer/solr-go.svg?style=shield)](https://circleci.com/gh/stevenferrer/solr-go)
[![Coverage Status](https://coveralls.io/repos/github/stevenferrer/solr-go/badge.svg?branch=master)](https://coveralls.io/github/stevenferrer/solr-go?branch=master)

# Solr-Go

[Solr](https://lucene.apache.org/solr/) client for [Go](http://go.dev/). 

```go

import (
    "context"
    // Import the package that you need
    solrquery "github.com/stevenferrer/solr-go/query"
)

func main() {
    // Initialize the query client
    queryClient := solrquery.NewClient("localhost", 8983)

    // Start querying!
    queryResp, err := queryClient.Query(
        context.Background(),
        "techproducts", // name of your collection
        map[string]interface{}{
            "query": "{!lucene df=name v=iPod}",
        },
    )
    ...
}
```

## Notes
* I'm using my project as the testbed for this module
* This is a *WORK IN-PROGRESS*, API might change a lot before *v1*
* Tested on [Solr 8.5](https://lucene.apache.org/solr/guide/8_5/)

## Contents

- [Solr-Go](#solr-go)
  - [Notes](#notes)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Features](#features)
  - [Contributing](#contributing)

## Installation

You can include it in your *go.mod* by running in your terminal (assuming you're inside the project directory):

```console
$ go get github.com/stevenferrer/solr-go
```

## Usage

Detailed documentation shall follow. For now you can start looking at the examples inside each package directory.

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
  - [x] JSON client
    - [x] [Add multiple documents](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-multiple-json-documents)
    - [x] [Multiple update commands](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands) 
	- [x] Example
  - [ ] XML client ??
  - [ ] CSV client ??
- [x] [JSON query API client](https://lucene.apache.org/solr/guide/8_5/json-query-dsl.html) client
  - [x] Facets
  - [x] Example
- [ ] [Standard query API client](https://lucene.apache.org/solr/guide/8_5/the-standard-query-parser.html#the-standard-query-parser) client ??
  - [ ] Example
- [x] Suggester client
- [x] Unified solr client
- [ ] Collections API client
- [ ] Configset API client
- [ ] Config API client
- [ ] Basic auth support
- [ ] Documentation

## Contributing

This is a **work in-progress**, any contributions are very welcome!

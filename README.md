[![Go Report Card](https://goreportcard.com/badge/github.com/stevenferrer/helios)](https://goreportcard.com/report/github.com/stevenferrer/helios)
[![CircleCI](https://circleci.com/gh/stevenferrer/helios.svg?style=shield)](https://circleci.com/gh/stevenferrer/helios)
[![Coverage Status](https://coveralls.io/repos/github/stevenferrer/helios/badge.svg?branch=master)](https://coveralls.io/github/stevenferrer/helios?branch=master)

# Helios

Helios contains set of packages for interacting with [Apache Solr](https://lucene.apache.org/solr/).

## Contents

- [Helios](#helios)
  - [Contents](#contents)
  - [Features and todo](#features-and-todo)
  - [Examples](#examples)
    - [Indexing and updating documents](#indexing-and-updating-documents)
	- [Quering with JSON request API](#querying-with-json-request-api)
	- [Interacting with Schema API](#interacting-with-schema-api)
	  - [Retrieving schema information](#retrieving-schema-information)
	  - [Managing and updating the schema](#modifying-and-updating-the-schema)
  - [Contributing](#contributing)

## Features and todo

- [x] [Schema API](https://lucene.apache.org/solr/guide/8_5/schema-api.html) client
  - [x] [Modifying schema](https://lucene.apache.org/solr/guide/8_5/schema-api.html#modify-the-schema)
  - [x] [Retrieving schema information](https://lucene.apache.org/solr/guide/8_5/schema-api.html#retrieve-schema-information)
  - [x] Example
- [ ] Index API
  - [x] JSON client
    - [x] [Add single document](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-a-single-json-document)
    - [x] [Add multiple documents](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#adding-multiple-json-documents)
    - [x] [Multiple update commands](https://lucene.apache.org/solr/guide/8_5/uploading-data-with-index-handlers.html#sending-json-update-commands) 
	- [x] Example
  - [ ] XML client ??
  - [ ] CSV client ??
- [x] [JSON Query API](https://lucene.apache.org/solr/guide/8_5/json-query-dsl.html) client
  - [x] Examples
- [ ] [Standard Query API](https://lucene.apache.org/solr/guide/8_5/the-standard-query-parser.html#the-standard-query-parser) client ??
  - [ ] Examples
- [ ] Basic Auth
- [ ] Improve documentation and create an awesome logo :)

## Examples

### Indexing and updating documents

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stevenferrer/helios"
	"github.com/stevenferrer/helios/index"
)

func main() {
    // Initialize index client
	host := "localhost"
	port := 8983
	client := index.NewJSONClient(host, port, &http.Client{
		Timeout: time.Second * 60,
	})

	var doc = struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		ID:   "1",
		Name: "Milana Vino",
	}

	collection := "gettingstarted"

	// Indexing a document
	err := client.AddSingle(context.Background(), collection, doc)
	if err != nil {
		log.Fatal(err)
	}

	// Indexing multiple documents
	var docs = []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		{
			ID:   "1",
			Name: "Milana Vino",
		},
		{
			ID:   "2",
			Name: "Charlie Jordan",
		},
		{
			ID:   "3",
			Name: "Daisy Keech",
		},
	}

	err = client.AddMultiple(context.Background(), collection, docs)
	if err != nil {
		log.Fatal(err)
	}

	// Sending multiple update commands
	err = client.UpdateCommands(context.Background(), collection,
		index.AddCommand{
			CommitWithin: 5000,
			Overwrite:    true,
			Doc: helios.M{
				"id":   "1",
				"name": "Milana Vino",
			},
		},
		index.AddCommand{
			Doc: helios.M{
				"id":   "2",
				"name": "Daisy Keech",
			},
		},
		index.AddCommand{
			Doc: helios.M{
				"id":   "3",
				"name": "Charley Jordan",
			},
		},
		index.DelByIDsCommand{
			IDs: []string{"2"},
		},
		index.DelByQryCommand{
			Query: "*:*",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

```

### Querying with JSON request API

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	. "github.com/stevenferrer/helios"
	"github.com/stevenferrer/helios/query"
)

func main() {
	// Initialize JSON query client
	host := "localhost"
	port := 8983
	queryClient := query.NewJSONClient(host, port, &http.Client{
		Timeout: time.Second * 60,
	})

	collection := "techproducts"

	// Simple query string
	resp, err := queryClient.Query(context.Background(), collection, M{
		"query": "name:iPod",
	})
	if err != nil {
		log.Fatal(err)
	}

	// do something with the response
	_ = resp

	// Full-expanded JSON object
	resp, err = queryClient.Query(context.Background(), collection, M{
		"query": M{
			"lucene": M{
				"df":    "name",
				"query": "iPod",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// do something with the response
	_ = resp

	// Complex queries
	resp, err = queryClient.Query(context.Background(), collection, M{
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

	if err != nil {
		log.Fatal(err)
	}

	// do something with the response
	_ = resp
}
```

### Interacting with Schema API

#### Retrieving schema information

```go

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stevenferrer/helios/schema"
)

func main() {
	host := "localhost"
	port := 8983

	// Initialize schema client
	schemaClient := schema.NewClient(host, port, &http.Client{
		Timeout: time.Second * 60,
	})

	collection := "gettingstarted"

	// Get the entire schema information
	gotSchema, err := schemaClient.GetSchema(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the value
	_ = gotSchema

	// List fields
	fields, err := schemaClient.ListFields(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the value
	_ = fields

	//* Get a specific field
	field, err := schemaClient.GetField(context.Background(), collection, "_text_")
	if err != nil {
		log.Fatal(err)
	}

	// do something with the value
	_ = field

	// List dynamic fields
	dynamicFields, err := schemaClient.ListDynamicFields(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}

	_ = dynamicFields

	// List field types
	fieldTypes, err := schemaClient.ListFieldTypes(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the value
	_ = fieldTypes

	// List copy fields
	copyFields, err := schemaClient.ListCopyFields(context.Background(), collection)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the value
	_ = copyFields
}
```

#### Modifying and updating the schema

```go
import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stevenferrer/helios/schema"
)

func main() {
	host := "localhost"
	port := 8983

	// Initialize schema client
	schemaClient := schema.NewClient(host, port, &http.Client{
		Timeout: time.Second * 60,
	})

	collection := "gettingstarted"

	// Adding a new field
	err := schemaClient.AddField(context.Background(), collection, schema.Field{
		Name:   "sell_by",
		Type:   "pdate",
		Stored: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Replacing a field
	err = schemaClient.ReplaceField(context.Background(), collection, schema.Field{
		Name:   "sell_by",
		Type:   "string",
		Stored: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Deleting a field
	err = schemaClient.DeleteField(context.Background(), collection, schema.Field{
		Name: "sell_by",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adding a dynamic field
	err = schemaClient.AddDynamicField(context.Background(), collection, schema.Field{
		Name:   "*_wtf",
		Type:   "string",
		Stored: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Replacing a dynamic field
	err = schemaClient.ReplaceDynamicField(context.Background(), collection, schema.Field{
		Name:   "*_wtf",
		Type:   "text_general",
		Stored: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Deleting a dynamic field
	err = schemaClient.DeleteDynamicField(context.Background(), collection, schema.Field{
		Name: "*_wtf",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adding a field type
	err = schemaClient.AddFieldType(context.Background(), collection, schema.FieldType{
		Name:  "myNewTextField",
		Class: "solr.TextField",
		IndexAnalyzier: &schema.Analyzer{
			Tokenizer: &schema.Tokenizer{
				Class:     "solr.PathHierarchyTokenizerFactory",
				Delimeter: "/",
			},
		},
		QueryAnalyzer: &schema.Analyzer{
			Tokenizer: &schema.Tokenizer{
				Class: "solr.KeywordTokenizerFactory",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Replacing a field type
	err = schemaClient.ReplaceFieldType(context.Background(), collection, schema.FieldType{
		Name:                 "myNewTextField",
		Class:                "solr.TextField",
		PositionIncrementGap: "100",
		Analyzer: &schema.Analyzer{
			Tokenizer: &schema.Tokenizer{
				Class: "solr.StandardTokenizerFactory",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Delete a field type
	err = schemaClient.DeleteFieldType(context.Background(), collection, schema.FieldType{
		Name: "myNewTextField",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adding a copy field
	err = schemaClient.AddCopyField(context.Background(), collection, schema.CopyField{
		Source: "*_shelf",
		Dest:   "_text_",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Deleting a copy field
	err = schemaClient.DeleteCopyField(context.Background(), collection, schema.CopyField{
		Source: "*_shelf",
		Dest:   "_text_",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

## Contributing

This is a work in-progres, any contributions are very welcome!

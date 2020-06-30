package main

import (
	"context"
	"log"

	"github.com/stevenferrer/solr-go/query"
)

type m = map[string]interface{}

func main() {
	// Initialize JSON query client
	host := "localhost"
	port := 8983
	queryClient := query.NewClient(host, port)

	collection := "techproducts"

	// Simple query string
	resp, err := queryClient.Query(context.Background(),
		collection, m{"query": "name:iPod"},
	)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the response
	_ = resp

	// Full-expanded JSON object
	resp, err = queryClient.Query(context.Background(),
		collection, m{
			"query": m{
				"lucene": m{
					"df":    "name",
					"query": "iPod",
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// do something with the response
	_ = resp

	// Complex queries
	resp, err = queryClient.Query(context.Background(), collection, m{
		"query": m{
			"boost": m{
				"query": m{
					"lucene": m{
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

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

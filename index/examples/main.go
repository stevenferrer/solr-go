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
	indexClient := index.NewJSONClient(host, port, &http.Client{
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
	err := indexClient.AddSingle(context.Background(), collection, doc)
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

	err = indexClient.AddMultiple(context.Background(), collection, docs)
	if err != nil {
		log.Fatal(err)
	}

	// Sending multiple update commands
	err = indexClient.UpdateCommands(context.Background(), collection,
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

package main

import (
	"context"
	"log"

	"github.com/stevenferrer/solr-go/index"
	"github.com/stevenferrer/solr-go/types"
)

func main() {
	// Initialize index client
	host := "localhost"
	port := 8983
	indexClient := index.NewJSONClient(host, port)

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
			Doc: types.M{
				"id":   "1",
				"name": "Milana Vino",
			},
		},
		index.AddCommand{
			Doc: types.M{
				"id":   "2",
				"name": "Daisy Keech",
			},
		},

		index.AddCommand{
			Doc: types.M{
				"id":   "3",
				"name": "Charlie Jordan",
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

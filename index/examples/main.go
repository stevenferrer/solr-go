package main

import (
	"context"
	"log"

	"github.com/stevenferrer/solr-go/index"
)

type m = map[string]interface{}

func main() {
	// Initialize index client
	host := "localhost"
	port := 8983
	indexClient := index.NewJSONClient(host, port)

	collection := "gettingstarted"

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
			Name: "Charly Jordan",
		},
		{
			ID:   "3",
			Name: "Daisy Keech",
		},
	}

	err := indexClient.AddDocs(context.Background(), collection, docs)
	if err != nil {
		log.Fatal(err)
	}

	// Sending multiple update commands
	err = indexClient.UpdateCommands(context.Background(), collection,
		index.AddCommand{
			CommitWithin: 5000,
			Overwrite:    true,
			Doc: m{
				"id":   "1",
				"name": "Milana Vino",
			},
		},
		index.AddCommand{
			Doc: m{
				"id":   "2",
				"name": "Daisy Keech",
			},
		},

		index.AddCommand{
			Doc: m{
				"id":   "3",
				"name": "Charly Jordan",
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

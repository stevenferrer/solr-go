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
	indexClient := index.NewClient(host, port)

	collection := "gettingstarted"

	// Indexing multiple documents
	var names = []struct {
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

	docs := index.NewDocs()
	for _, name := range names {
		docs.AddDoc(name)
	}

	ctx := context.Background()
	err := indexClient.AddDocuments(ctx, collection, docs)
	checkErr(err)
	err = indexClient.Commit(ctx, collection)
	checkErr(err)

	// Sending multiple update commands
	err = indexClient.SendCommands(context.Background(), collection,
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
		index.DeleteByIDsCommand{
			IDs: []string{"2"},
		},
		index.DeleteByQueryCommand{
			Query: "*:*",
		},
	)
	checkErr(err)
	err = indexClient.Commit(ctx, collection)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

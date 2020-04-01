package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/stevenferrer/helios"
)

func main() {
	ctx := context.Background()

	const (
		host = "192.168.100.19"
		port = 8983
	)

	client := helios.NewClient(host, port)
	b, err := json.Marshal(helios.SimpleQueryRequest{
		Query: "*:*",
		Facet: helios.M{
			"categories": helios.M{
				"type":  "terms",
				"field": "category",
				"limit": 100,
			},
		},
	})
	checkErr(err)

	fmt.Println(string(b))

	response, err := client.Query(ctx, "products", b)
	checkErr(err)

	b, err = json.Marshal(response)
	checkErr(err)

	err = ioutil.WriteFile("response.json", b, 0644)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

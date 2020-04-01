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

	// solr server configurations
	const (
		host = "192.168.100.19"
		port = 8983
	)

	// create a new helios client
	client := helios.NewClient(host, port)
	b, err := ioutil.ReadFile("beers.json")
	checkErr(err)

	uresponse, err := client.Update(ctx, "products", b)
	checkErr(err)

	b, err = json.Marshal(uresponse)
	fmt.Println("update resposne")
	fmt.Println(string(b))

	b, err = json.Marshal(helios.SimpleSelectRequest{
		Query: "*:*",
	})
	checkErr(err)

	// perform query
	sresponse, err := client.Select(ctx, "products", b)
	checkErr(err)

	b, err = json.MarshalIndent(sresponse, "", "  ")
	checkErr(err)

	// print response
	fmt.Println("select reponse")
	fmt.Println(string(b))

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

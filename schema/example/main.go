package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/stevenferrer/helios/schema"
)

func retrievingSchema() {
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

func modifyingSchema() {
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

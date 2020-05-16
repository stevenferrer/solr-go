package schema_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/helios/schema"
)

func TestRetrieveSchema(t *testing.T) {
	ctx := context.Background()
	collection := "gettingstarted"
	host := "localhost"
	port := 8983

	t.Run("get schema", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/get-schema")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			gotSchema, err := client.GetSchema(ctx, collection)
			a.NoError(err)
			a.NotNil(gotSchema)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.GetSchema(ctx, collection)
				assert.Error(t, err)
			})
		})
	})

	t.Run("list fields", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/list-fields")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			gotFields, err := client.ListFields(ctx, collection)
			a.NoError(err)
			a.NotNil(gotFields)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.ListFields(ctx, collection)
				assert.Error(t, err)
			})
		})
	})

	t.Run("get field", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/get-field")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			gotField, err := client.GetField(ctx, collection, "_text_")
			a.NoError(err)
			a.NotNil(gotField)
		})
		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.GetField(ctx, collection, "_text_")
				assert.Error(t, err)
			})
		})
	})

	t.Run("list dynamic fields", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/list-dynamic-fields")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			gotDynamicFields, err := client.ListDynamicFields(ctx, collection)
			a.NoError(err)
			a.NotNil(gotDynamicFields)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.ListDynamicFields(ctx, collection)
				assert.Error(t, err)
			})
		})
	})

	t.Run("list field types", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/list-field-types")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			gotFieldTypes, err := client.ListFieldTypes(ctx, collection)
			a.NoError(err)
			a.NotNil(gotFieldTypes)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.ListFieldTypes(ctx, collection)
				assert.Error(t, err)
			})
		})
	})

	t.Run("list copy fields", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			a := assert.New(t)

			rec, err := recorder.New("fixtures/list-copy-fields")
			require.NoError(t, err)
			defer rec.Stop()

			client := schema.NewClient(host, port, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			err = client.AddField(ctx, collection, schema.Field{
				Name:   "my_field",
				Type:   "string",
				Stored: true,
			})
			require.NoError(t, err)

			err = client.AddCopyField(ctx, collection, schema.CopyField{
				Source: "my_field",
				Dest:   "_text_",
			})
			require.NoError(t, err)

			var gotCopyFields []schema.CopyField
			gotCopyFields, err = client.ListCopyFields(ctx, collection)
			a.NoError(err)
			a.NotNil(gotCopyFields)
		})

		t.Run("error", func(t *testing.T) {
			t.Run("parse url", func(t *testing.T) {
				client := schema.NewClient("wtf\\:\\wtf", port, &http.Client{})

				_, err := client.ListCopyFields(ctx, collection)
				assert.Error(t, err)
			})
		})
	})
}

func TestModifySchema(t *testing.T) {
	var testCases = []struct {
		name, command string
		body          interface{}
		wantErr       bool
	}{
		{
			name:    "add field with error",
			command: "add-field",
			body:    schema.Field{},
			wantErr: true,
		},

		{
			name:    "add field ok",
			command: "add-field",
			body: schema.Field{
				Name:   "sell_by",
				Type:   "pdate",
				Stored: true,
			},
		},

		{
			name:    "replace field with error",
			command: "replace-field",
			body:    schema.Field{},
			wantErr: true,
		},
		{
			name:    "replace field ok",
			command: "replace-field",
			body: schema.Field{
				Name:   "sell_by",
				Type:   "string",
				Stored: false,
			},
		},

		{
			name:    "delete field with error",
			command: "delete-field",
			body:    schema.Field{},
			wantErr: true,
		},
		{
			name:    "delete field ok",
			command: "delete-field",
			body: schema.Field{
				Name: "sell_by",
			},
		},

		{
			name:    "add dynamic field with error",
			command: "add-dynamic-field",
			body:    schema.Field{},
			wantErr: true,
		},
		{
			name:    "add dynamic field ok",
			command: "add-dynamic-field",
			body: schema.Field{
				Name:   "*_wtf",
				Type:   "string",
				Stored: true,
			},
		},
		{
			name:    "replace dynamic field with error",
			command: "replace-dynamic-field",
			body:    schema.Field{},
			wantErr: true,
		},
		{
			name:    "replace dynamic field ok",
			command: "replace-dynamic-field",
			body: schema.Field{
				Name:   "*_wtf",
				Type:   "text_general",
				Stored: false,
			},
		},
		{
			name:    "delete dynamic field with error",
			command: "delete-dynamic-field",
			body:    schema.Field{},
			wantErr: true,
		},
		{
			name:    "delete dynamic field ok",
			command: "delete-dynamic-field",
			body:    schema.Field{Name: "*_wtf"},
		},
		{
			name:    "add field type with error",
			command: "add-field-type",
			body:    schema.FieldType{},
			wantErr: true,
		},
		{
			name:    "add field type ok",
			command: "add-field-type",
			body: schema.FieldType{
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
			},
		},

		{
			name:    "replace field type with error",
			command: "replace-field-type",
			body:    schema.FieldType{},
			wantErr: true,
		},
		{
			name:    "replace field type ok",
			command: "replace-field-type",
			body: schema.FieldType{
				Name:                 "myNewTextField",
				Class:                "solr.TextField",
				PositionIncrementGap: "100",
				Analyzer: &schema.Analyzer{
					Tokenizer: &schema.Tokenizer{
						Class: "solr.StandardTokenizerFactory",
					},
				},
			},
		},

		{
			name:    "delete field type with error",
			command: "delete-field-type",
			body:    schema.FieldType{},
			wantErr: true,
		},
		{
			name:    "delete field type ok",
			command: "delete-field-type",
			body: schema.FieldType{
				Name: "myNewTextField",
			},
		},

		{
			name:    "add copy field with error",
			command: "add-copy-field",
			body:    schema.CopyField{},
			wantErr: true,
		},
		{
			name:    "add copy field ok",
			command: "add-copy-field",
			body: schema.CopyField{
				Source: "*_shelf",
				Dest:   "_text_",
			},
		},

		{
			name:    "delete copy field with error",
			command: "delete-copy-field",
			body:    schema.CopyField{},
			wantErr: true,
		},
		{
			name:    "delete copy field ok",
			command: "delete-copy-field",
			body: schema.CopyField{
				Source: "*_shelf",
				Dest:   "_text_",
			},
		},
	}

	ctx := context.Background()
	coll := "gettingstarted"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := fmt.Sprintf("fixtures/%s", strings.ReplaceAll(tc.name, " ", "-"))

			r, err := recorder.New(fixture)
			require.NoError(t, err)
			defer r.Stop()

			httpClient := &http.Client{
				Timeout:   time.Second * 60,
				Transport: r,
			}

			client := schema.NewClient("localhost", 8983, httpClient)

			switch tc.command {
			case "add-field":
				err = client.AddField(ctx, coll, tc.body.(schema.Field))
			case "delete-field":
				err = client.DeleteField(ctx, coll, tc.body.(schema.Field))
			case "replace-field":
				err = client.ReplaceField(ctx, coll, tc.body.(schema.Field))
			case "add-dynamic-field":
				err = client.AddDynamicField(ctx, coll, tc.body.(schema.Field))
			case "delete-dynamic-field":
				err = client.DeleteDynamicField(ctx, coll, tc.body.(schema.Field))
			case "replace-dynamic-field":
				err = client.ReplaceDynamicField(ctx, coll, tc.body.(schema.Field))
			case "add-field-type":
				err = client.AddFieldType(ctx, coll, tc.body.(schema.FieldType))
			case "delete-field-type":
				err = client.DeleteFieldType(ctx, coll, tc.body.(schema.FieldType))
			case "replace-field-type":
				err = client.ReplaceFieldType(ctx, coll, tc.body.(schema.FieldType))
			case "add-copy-field":
				cpyField := tc.body.(schema.CopyField)
				if !tc.wantErr {
					err = client.AddDynamicField(ctx, coll, schema.Field{
						Name:   cpyField.Source,
						Type:   "string",
						Stored: true,
					})
					assert.NoError(t, err)
				}
				err = client.AddCopyField(ctx, coll, cpyField)
			case "delete-copy-field":
				err = client.DeleteCopyField(ctx, coll, tc.body.(schema.CopyField))
			default:
				err = errors.New("command not found")
			}

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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

	"github.com/stevenferrer/helios"
	"github.com/stevenferrer/helios/schema"
)

func TestClient(t *testing.T) {
	var testCases = []struct {
		name, command string
		m             helios.M
		wantErr       bool
	}{
		{
			name:    "add field with error",
			command: "add-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "add field ok",
			command: "add-field",
			m: helios.M{
				"name":   "sell_by",
				"type":   "pdate",
				"stored": true,
			},
		},

		{
			name:    "replace field with error",
			command: "replace-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "replace field ok",
			command: "replace-field",
			m: helios.M{
				"name":   "sell_by",
				"type":   "string",
				"stored": false,
			},
		},

		{
			name:    "delete field with error",
			command: "delete-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "delete field ok",
			command: "delete-field",
			m: helios.M{
				"name": "sell_by",
			},
		},

		{
			name:    "add dynamic field with error",
			command: "add-dynamic-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "add dynamic field ok",
			command: "add-dynamic-field",
			m: helios.M{
				"name":   "*_wtf",
				"type":   "string",
				"stored": true,
			},
		},
		{
			name:    "replace dynamic field with error",
			command: "replace-dynamic-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "replace dynamic field ok",
			command: "replace-dynamic-field",
			m: helios.M{
				"name":   "*_wtf",
				"type":   "text_general",
				"stored": false,
			},
		},
		{
			name:    "delete dynamic field with error",
			command: "delete-dynamic-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "delete dynamic field ok",
			command: "delete-dynamic-field",
			m:       helios.M{"name": "*_wtf"},
		},
		{
			name:    "add field type with error",
			command: "add-field-type",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "add field type ok",
			command: "add-field-type",
			m: helios.M{
				"name":  "myNewTextField",
				"class": "solr.TextField",
				"indexAnalyzer": helios.M{
					"tokenizer": helios.M{
						"class":     "solr.PathHierarchyTokenizerFactory",
						"delimiter": "/",
					},
				},
				"queryAnalyzer": helios.M{
					"tokenizer": helios.M{
						"class": "solr.KeywordTokenizerFactory",
					},
				},
			},
		},

		{
			name:    "replace field type with error",
			command: "replace-field-type",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "replace field type ok",
			command: "replace-field-type",
			m: helios.M{
				"name":                 "myNewTextField",
				"class":                "solr.TextField",
				"positionIncrementGap": "100",
				"analyzer": helios.M{
					"tokenizer": helios.M{
						"class": "solr.StandardTokenizerFactory",
					},
				},
			},
		},

		{
			name:    "delete field type with error",
			command: "delete-field-type",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "delete field type ok",
			command: "delete-field-type",
			m:       helios.M{"name": "myNewTextField"},
		},

		{
			name:    "add copy field with error",
			command: "add-copy-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "add copy field ok",
			command: "add-copy-field",
			m: helios.M{
				"source": "*_shelf",
				"dest":   "_text_",
			},
		},

		{
			name:    "delete copy field with error",
			command: "delete-copy-field",
			m:       helios.M{},
			wantErr: true,
		},
		{
			name:    "delete copy field ok",
			command: "delete-copy-field",
			m: helios.M{
				"source": "*_shelf",
				"dest":   "_text_",
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
				err = client.AddField(ctx, coll, tc.m)
			case "delete-field":
				err = client.DeleteField(ctx, coll, tc.m)
			case "replace-field":
				err = client.ReplaceField(ctx, coll, tc.m)
			case "add-dynamic-field":
				err = client.AddDynamicField(ctx, coll, tc.m)
			case "delete-dynamic-field":
				err = client.DeleteDynamicField(ctx, coll, tc.m)
			case "replace-dynamic-field":
				err = client.ReplaceDynamicField(ctx, coll, tc.m)
			case "add-field-type":
				err = client.AddFieldType(ctx, coll, tc.m)
			case "delete-field-type":
				err = client.DeleteFieldType(ctx, coll, tc.m)
			case "replace-field-type":
				err = client.ReplaceFieldType(ctx, coll, tc.m)
			case "add-copy-field":
				if !tc.wantErr {
					err = client.AddDynamicField(ctx, coll, helios.M{
						"name":   tc.m["source"],
						"type":   "string",
						"stored": true,
					})
					assert.NoError(t, err)
				}
				err = client.AddCopyField(ctx, coll, tc.m)
			case "delete-copy-field":
				err = client.DeleteCopyField(ctx, coll, tc.m)
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

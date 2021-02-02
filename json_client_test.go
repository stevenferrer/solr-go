package solr_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/solr-go"
)

func TestJSONClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	baseURL := "https://solr.example.com"
	collection := "products"

	client := solr.NewJSONClient(baseURL).
		WithHTTPClient(&http.Client{
			Timeout: 10 * time.Second,
		})

	t.Run("query", func(t *testing.T) {
		mockBody := `{"query":"{!dismax}apple pie"}`
		httpmock.RegisterResponder(
			http.MethodPost,
			baseURL+"/solr/"+collection+"/query",
			newResponder(mockBody, solr.M{}),
		)

		q := solr.NewQuery().QueryParser(
			solr.NewDisMaxQueryParser("apple pie"))

		_, err := client.Query(ctx, collection, q)
		assert.NoError(t, err)
	})

	t.Run("index", func(t *testing.T) {
		mockBody := `[{"id":1,"name":"product 1"},{"id":2,"name":"product 2"},{"id":3,"name":"product 3"}]`
		httpmock.RegisterResponder(
			http.MethodPost,
			baseURL+"/solr/"+collection+"/update",
			newResponder(mockBody, solr.M{}),
		)

		httpmock.RegisterResponder(
			http.MethodGet,
			baseURL+"/solr/"+collection+"/update",
			func(r *http.Request) (*http.Response, error) {
				commit := r.URL.Query().Get("commit")
				if commit != "true" {
					return nil, errors.New("`commit` param to be true")
				}

				return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
			},
		)

		doc1 := solr.Document{
			"id":   1,
			"name": "product 1",
		}

		doc2 := solr.Document{
			"id":   2,
			"name": "product 2",
		}
		doc3 := solr.Document{
			"id":   3,
			"name": "product 3",
		}

		_, err := client.Update(ctx, collection, doc1, doc2, doc3)
		assert.NoError(t, err)

		err = client.Commit(ctx, collection)
		assert.NoError(t, err)
	})

	t.Run("schema", func(t *testing.T) {
		t.Run("add fields", func(t *testing.T) {
			mockBody := `{"add-field":[{"name":"foo","type":"string"},{"name":"bar","type":"string"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name: "foo",
					Type: "string",
				},
				{
					Name: "bar",
					Type: "string",
				},
			}
			err := client.AddFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("delete fields", func(t *testing.T) {
			mockBody := `{"delete-field":[{"name":"foo"},{"name":"bar"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			}
			err := client.DeleteFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("replace fields", func(t *testing.T) {
			mockBody := `{"replace-field":[{"name":"foo","type":"plong"},{"name":"bar","type":"plong"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name: "foo",
					Type: "plong",
				},
				{
					Name: "bar",
					Type: "plong",
				},
			}
			err := client.ReplaceFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("add dynamic fields", func(t *testing.T) {
			mockBody := `{"add-dynamic-field":[{"name":"*_foo","type":"string","stored":true},{"name":"*_bar","type":"plong"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name:   "*_foo",
					Type:   "string",
					Stored: true,
				},
				{
					Name:   "*_bar",
					Type:   "plong",
					Stored: false,
				},
			}
			err := client.AddDynamicFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("delete dynamic fields", func(t *testing.T) {
			mockBody := `{"delete-dynamic-field":[{"name":"*_foo"},{"name":"*_bar"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name: "*_foo",
				},
				{
					Name: "*_bar",
				},
			}
			err := client.DeleteDynamicFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("replace dynamic fields", func(t *testing.T) {
			mockBody := `{"replace-dynamic-field":[{"name":"*_foo","type":"text_general"},{"name":"*_bar","type":"string"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fields := []solr.Field{
				{
					Name: "*_foo",
					Type: "text_general",
				},
				{
					Name: "*_bar",
					Type: "string",
				},
			}
			err := client.ReplaceDynamicFields(ctx, collection, fields...)
			require.NoError(t, err)
		})

		t.Run("add field types", func(t *testing.T) {
			mockBody := `{"add-field-type":[{"name":"myNewTextField","class":"solr.TextField","indexAnalyzer":{"tokenizer":{"class":"solr.PathHierarchyTokenizerFactory","delimiter":"/"}},"queryAnalyzer":{"tokenizer":{"class":"solr.KeywordTokenizerFactory"}}}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fieldTypes := []solr.FieldType{
				{
					Name:  "myNewTextField",
					Class: "solr.TextField",
					IndexAnalyzer: &solr.Analyzer{
						Tokenizer: &solr.Tokenizer{
							Class:     "solr.PathHierarchyTokenizerFactory",
							Delimeter: "/",
						},
					},
					QueryAnalyzer: &solr.Analyzer{
						Tokenizer: &solr.Tokenizer{
							Class: "solr.KeywordTokenizerFactory",
						},
					},
				},
			}
			err := client.AddFieldTypes(ctx, collection, fieldTypes...)
			require.NoError(t, err)
		})

		t.Run("delete field types", func(t *testing.T) {
			mockBody := `{"delete-field-type":[{"name":"myNewTextField"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fieldTypes := []solr.FieldType{
				{
					Name: "myNewTextField",
				},
			}
			err := client.DeleteFieldTypes(ctx, collection, fieldTypes...)
			require.NoError(t, err)
		})

		t.Run("replace field types", func(t *testing.T) {
			mockBody := `{"replace-field-type":[{"name":"myNewTextField","class":"solr.TextField","positionIncrementGap":"100","analyzer":{"tokenizer":{"class":"solr.StandardTokenizerFactory"}}}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			fieldTypes := []solr.FieldType{
				{
					Name:                 "myNewTextField",
					Class:                "solr.TextField",
					PositionIncrementGap: "100",
					Analyzer: &solr.Analyzer{
						Tokenizer: &solr.Tokenizer{
							Class: "solr.StandardTokenizerFactory",
						},
					},
				},
			}
			err := client.ReplaceFieldTypes(ctx, collection, fieldTypes...)
			require.NoError(t, err)
		})

		t.Run("add copy fields", func(t *testing.T) {
			mockBody := `{"add-copy-field":[{"source":"shelf","dest":"location"},{"source":"shelf","dest":"catchall"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			copyFields := []solr.CopyField{
				{
					Source: "shelf",
					Dest:   "location",
				},
				{
					Source: "shelf",
					Dest:   "catchall",
				},
			}
			err := client.AddCopyFields(ctx, collection, copyFields...)
			require.NoError(t, err)
		})

		t.Run("delete copy fields", func(t *testing.T) {
			mockBody := `{"delete-copy-field":[{"source":"shelf","dest":"location"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponder(mockBody, solr.M{}),
			)

			copyFields := []solr.CopyField{
				{
					Source: "shelf",
					Dest:   "location",
				},
			}
			err := client.DeleteCopyFields(ctx, collection, copyFields...)
			require.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			mockBody := `{"add-field":[{"name":"foo"}]}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/schema",
				newResponderWithStatus(http.StatusBadRequest, mockBody, solr.BaseResponse{
					Error: &solr.Error{
						Msg: "this is an error",
					},
				}),
			)

			fields := []solr.Field{{Name: "foo"}}
			err := client.AddFields(ctx, collection, fields...)
			assert.Error(t, err)
		})
	})

	t.Run("config", func(t *testing.T) {
		t.Run("set property", func(t *testing.T) {
			mockBody := `{"set-property":{"updater.autoCommit.maxTime":15000}}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			err := client.SetProperties(ctx, collection, solr.CommonProperty{
				Name:  "updater.autoCommit.maxTime",
				Value: 15000,
			})
			assert.NoError(t, err)
		})

		t.Run("unset property", func(t *testing.T) {
			mockBody := `{"unset-property":"updater.autoCommit.maxTime"}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			err := client.UnsetProperty(ctx, collection, solr.CommonProperty{
				Name: "updater.autoCommit.maxTime",
			})
			assert.NoError(t, err)
		})

		t.Run("add component", func(t *testing.T) {
			mockBody := `{"add-requesthandler":{"class":"solr.DumpRequestHandler","defaults":{"a":"b","rows":10,"x":"y"},"name":"/mypath","useParams":"x"}}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			err := client.AddComponent(ctx, collection, solr.Component{
				Type:  solr.RequestHandler,
				Name:  "/mypath",
				Class: "solr.DumpRequestHandler",
				Values: solr.M{
					"defaults":  solr.M{"x": "y", "a": "b", "rows": 10},
					"useParams": "x",
				},
			})
			assert.NoError(t, err)
		})

		t.Run("update component", func(t *testing.T) {
			mockBody := `{"update-requesthandler":{"class":"solr.DumpRequestHandler","defaults":{"rows":"20","x":"new value for x"},"name":"/mypath","useParams":"x"}}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			err := client.UpdateComponent(ctx, collection, solr.Component{
				Type:  solr.RequestHandler,
				Name:  "/mypath",
				Class: "solr.DumpRequestHandler",
				Values: solr.M{
					"defaults":  solr.M{"x": "new value for x", "rows": "20"},
					"useParams": "x",
				},
			})
			assert.NoError(t, err)
		})

		t.Run("delete component", func(t *testing.T) {
			mockBody := `{"delete-requesthandler":"/mypath"}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			err := client.DeleteComponent(ctx, collection, solr.Component{
				Type: solr.RequestHandler,
				Name: "/mypath",
			})
			assert.NoError(t, err)
		})
	})

	t.Run("suggest", func(t *testing.T) {
		responder, err := httpmock.NewJsonResponder(http.StatusOK, solr.SuggestResponse{})
		require.NoError(t, err)
		httpmock.RegisterResponder(
			http.MethodGet,
			baseURL+"/solr/"+collection+"/suggest?suggest=true&suggest.build=true&suggest.dictionary=mySuggester&suggest.q=elec",
			responder,
		)

		suggestParams := solr.NewSuggesterParams("suggest").
			Build().Dictionaries("mySuggester").Query("elec")
		_, err = client.Suggest(ctx, collection, suggestParams)
		assert.NoError(t, err)
	})
}

func newResponder(body string, mockResp interface{}) httpmock.Responder {
	return newResponderWithStatus(http.StatusOK, body, mockResp)
}

func newResponderWithStatus(status int, body string, mockResp interface{}) httpmock.Responder {
	var mockBody interface{}
	err := json.Unmarshal([]byte(body), &mockBody)
	if err != nil {
		panic(err)
	}

	return func(r *http.Request) (*http.Response, error) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, "read request body")
		}

		var reqBody interface{}
		err = json.Unmarshal(b, &reqBody)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal request body")
		}

		if !reflect.DeepEqual(reqBody, mockBody) {
			return nil, errors.Errorf("expected request body: %v", string(b))
		}

		return httpmock.NewJsonResponse(status, mockResp)
	}
}

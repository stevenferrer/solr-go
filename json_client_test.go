package solr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/solr-go"
)

// errorRequestSender is a request sender that errors
type errorRequestSender struct{}

var _ solr.RequestSender = (*errorRequestSender)(nil)

var errSendRequest = errors.New("an error from request sender")

func (rs *errorRequestSender) SendRequest(_ context.Context, _, _, _ string, _ io.Reader) (*http.Response, error) {
	return nil, errSendRequest
}

func TestJSONClientMock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	baseURL := "https://solr.example.com"
	collection := "products"

	client := solr.NewJSONClient(baseURL)
	clientThatErrors := solr.NewJSONClient(baseURL).
		WithRequestSender(&errorRequestSender{})

	t.Run("collections", func(t *testing.T) {
		t.Run("create collection", func(t *testing.T) {
			httpmock.RegisterResponder(
				http.MethodGet,
				baseURL+"/solr/admin/collections",
				func(r *http.Request) (*http.Response, error) {
					query := "action=CREATE&name=mycollection&numShards=1&replicationFactor=1"
					gotQuery := r.URL.Query().Encode()
					if gotQuery != query {
						return nil, errors.Errorf("expecting url query to be %q but got %q", query, gotQuery)
					}

					return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
				},
			)

			params := solr.NewCollectionParams().
				Name("mycollection").NumShards(1).
				ReplicationFactor(1)
			err := client.CreateCollection(ctx, params)
			assert.NoError(t, err)

			err = clientThatErrors.CreateCollection(ctx, params)
			assert.ErrorIs(t, err, errSendRequest)
		})
		t.Run("delete collection", func(t *testing.T) {
			httpmock.RegisterResponder(
				http.MethodGet,
				baseURL+"/solr/admin/collections",
				func(r *http.Request) (*http.Response, error) {
					query := "action=DELETE&name=mycollection"
					gotQuery := r.URL.Query().Encode()
					if gotQuery != query {
						return nil, errors.Errorf("expecting url query to be %q but got %q", query, gotQuery)
					}

					return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
				},
			)

			params := solr.NewCollectionParams().
				Name("mycollection")
			err := client.DeleteCollection(ctx, params)
			assert.NoError(t, err)

			err = clientThatErrors.DeleteCollection(ctx, params)
			assert.ErrorIs(t, err, errSendRequest)
		})
	})

	t.Run("core admin", func(t *testing.T) {
		t.Run("create core", func(t *testing.T) {
			httpmock.RegisterResponder(
				http.MethodGet,
				baseURL+"/solr/admin/cores",
				func(r *http.Request) (*http.Response, error) {
					query := "action=CREATE&name=mycore"
					gotQuery := r.URL.Query().Encode()
					if gotQuery != query {
						return nil, errors.Errorf("expecting url query to be %q but got %q", query, gotQuery)
					}

					return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
				},
			)

			params := solr.NewCreateCoreParams("mycore")
			err := client.CreateCore(ctx, params)
			assert.NoError(t, err)

			err = clientThatErrors.CreateCore(ctx, params)
			assert.ErrorIs(t, err, errSendRequest)
		})

		t.Run("core status", func(t *testing.T) {
			httpmock.RegisterResponder(
				http.MethodGet,
				baseURL+"/solr/admin/cores",
				func(r *http.Request) (*http.Response, error) {
					query := "action=STATUS&core=mycore"
					gotQuery := r.URL.Query().Encode()
					if gotQuery != query {
						return nil, errors.Errorf("expecting url query to be %q but got %q", query, gotQuery)
					}

					return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
				},
			)

			params := solr.NewCoreParams("mycore")
			_, err := client.CoreStatus(ctx, params)
			assert.NoError(t, err)

			_, err = clientThatErrors.CoreStatus(ctx, params)
			assert.ErrorIs(t, err, errSendRequest)

		})

		t.Run("unload core", func(t *testing.T) {
			httpmock.RegisterResponder(
				http.MethodGet,
				baseURL+"/solr/admin/cores",
				func(r *http.Request) (*http.Response, error) {
					query := "action=UNLOAD&core=mycore"
					gotQuery := r.URL.Query().Encode()
					if gotQuery != query {
						return nil, errors.Errorf("expecting url query to be %q but got %q", query, gotQuery)
					}

					return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
				},
			)

			params := solr.NewCoreParams("mycore")
			err := client.UnloadCore(ctx, params)
			assert.NoError(t, err)

		})
	})

	t.Run("query", func(t *testing.T) {
		mockBody := `{"query":"{!dismax v='apple pie'}"}`
		httpmock.RegisterResponder(
			http.MethodPost,
			baseURL+"/solr/"+collection+"/query",
			newResponder(mockBody, solr.M{}),
		)

		query := solr.NewQuery(solr.NewDisMaxQueryParser().
			Query("'apple pie'").BuildParser())
		_, err := client.Query(ctx, collection, query)
		assert.NoError(t, err)

		_, err = clientThatErrors.Query(ctx, collection, query)
		assert.ErrorIs(t, err, errSendRequest)
	})

	t.Run("update and commit", func(t *testing.T) {
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
					return nil, errors.New("expect `commit` param to be true")
				}

				return httpmock.NewJsonResponse(http.StatusOK, solr.M{})
			},
		)

		buf := &bytes.Buffer{}
		err := json.NewEncoder(buf).
			Encode([]solr.M{
				{
					"id":   1,
					"name": "product 1",
				},
				{
					"id":   2,
					"name": "product 2",
				},
				{
					"id":   3,
					"name": "product 3",
				},
			})
		assert.NoError(t, err)

		_, err = client.Update(ctx, collection, solr.JSON, buf)
		assert.NoError(t, err)

		err = client.Commit(ctx, collection)
		assert.NoError(t, err)

		_, err = clientThatErrors.Update(ctx, collection, solr.JSON, buf)
		assert.ErrorIs(t, err, errSendRequest)

		err = clientThatErrors.Commit(ctx, collection)
		assert.ErrorIs(t, err, errSendRequest)
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

			err = clientThatErrors.AddFields(ctx, collection, fields...)
			assert.ErrorIs(t, err, errSendRequest)
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
					Error: &solr.ResponseError{
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

		t.Run("add components", func(t *testing.T) {
			mockBody := `{"add-searchcomponent":{"class":"solr.SuggestComponent","name":"suggest","suggester":{"dictionaryImpl":"DocumentDictionaryFactory","field":"suggest","lookupImpl":"AnalyzingInfixLookupFactory","name":"default","suggestAnalyzerFieldType":"suggest_text"}},"add-requesthandler":{"class":"solr.SearchHandler","components":["suggest"],"defaults":{"suggest":true,"suggest.count":10,"suggest.dictionary":"default"},"name":"/suggest","startup":"lazy"}}`
			httpmock.RegisterResponder(
				http.MethodPost,
				baseURL+"/solr/"+collection+"/config",
				newResponder(mockBody, solr.M{}),
			)

			suggestComponent := solr.NewComponent(solr.SearchComponent).
				Name("suggest").Class("solr.SuggestComponent").
				Config(solr.M{
					"suggester": solr.M{
						"name":                     "default",
						"lookupImpl":               "AnalyzingInfixLookupFactory",
						"dictionaryImpl":           "DocumentDictionaryFactory",
						"field":                    "suggest",
						"suggestAnalyzerFieldType": "suggest_text",
					},
				})

			suggestHandler := solr.NewComponent(solr.RequestHandler).
				Name("/suggest").Class("solr.SearchHandler").
				Config(solr.M{
					"startup": "lazy",
					"defaults": solr.M{
						"suggest":            true,
						"suggest.count":      10,
						"suggest.dictionary": "default",
					},
					"components": []string{"suggest"},
				})

			err := client.AddComponents(ctx, collection, suggestComponent, suggestHandler)
			assert.NoError(t, err)

			err = clientThatErrors.AddComponents(ctx, collection, suggestComponent)
			assert.ErrorIs(t, err, errSendRequest)
		})

		t.Run("update components", func(t *testing.T) {
			err := client.UpdateComponents(ctx, collection)
			assert.Error(t, err)
		})

		t.Run("delete components", func(t *testing.T) {
			err := client.DeleteComponents(ctx, collection)
			assert.Error(t, err)
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

		_, err = clientThatErrors.Suggest(ctx, collection, suggestParams)
		assert.ErrorIs(t, err, errSendRequest)
	})

	t.Run("unexpected html", func(t *testing.T) {
		httpmock.RegisterResponder(http.MethodGet, baseURL+"/solr/admin/cores", func(r *http.Request) (*http.Response, error) {
			response := httpmock.NewBytesResponse(http.StatusUnauthorized, []byte("<html><title>Unauthorized</html>"))
			response.Header.Set("Content-Type", "text/html")
			return response, nil
		})

		params := solr.NewCoreParams("mycore")
		_, err := client.CoreStatus(ctx, params)
		assert.Error(t, err)
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
		b, err := io.ReadAll(r.Body)
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

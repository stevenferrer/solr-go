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

	"github.com/sf9v/solr-go"
)

func TestJSONClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctx := context.Background()
	baseURL := "https://solr.example.com"
	collection := "products"

	t.Run("query", func(t *testing.T) {
		mockBody := `{"query":"{!dismax}apple pie"}`
		httpmock.RegisterResponder(
			http.MethodPost,
			baseURL+"/solr/"+collection+"/query",
			newMatchResponder(mockBody, solr.M{}),
		)

		client := solr.NewJSONClient(baseURL).
			WithHTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			})
		q := solr.NewQuery().WithQueryParser(
			solr.NewDisMaxQueryParser("apple pie"),
		)

		_, err := client.Query(ctx, collection, q)
		assert.NoError(t, err)
	})

	t.Run("update and commit", func(t *testing.T) {
		mockBody := `[{"id":1,"name":"product 1"},{"id":2,"name":"product 2"},{"id":3,"name":"product 3"}]`
		httpmock.RegisterResponder(
			http.MethodPost,
			baseURL+"/solr/"+collection+"/update",
			newMatchResponder(mockBody, solr.M{}),
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

		client := solr.NewJSONClient(baseURL).
			WithHTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			})

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
}

func newMatchResponder(body string, mockResp interface{}) httpmock.Responder {
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

		return httpmock.NewJsonResponse(http.StatusOK, mockResp)
	}
}

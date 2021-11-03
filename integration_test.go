//go:build integration
// +build integration

package solr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/solr-go"
)

func TestJSONClient(t *testing.T) {
	requestSender := solr.NewDefaultRequestSender().
		WithHTTPClient(&http.Client{Timeout: 30 * time.Second}).
		WithBasicAuth("solr", "SolrRocks")

	t.Run("standalone", func(t *testing.T) {
		baseURL := "http://localhost:8983"
		core := "searchengines"
		client := solr.NewJSONClient(baseURL).
			WithRequestSender(requestSender)

		// check core status
		ctx := context.Background()
		_, err := client.CoreStatus(ctx, solr.NewCoreParams(core))
		require.NoError(t, err, "core status should not error")

		newTestRunner(client, core)(t)

		err = client.UnloadCore(ctx, solr.NewCoreParams(core).
			DeleteIndex(true).DeleteInstanceDir(true).DeleteDataDir(true))
		require.NoError(t, err)
	})

	t.Run("solrcloud", func(t *testing.T) {
		baseURL := "http://localhost:8984"
		collection := "searchengines"
		client := solr.NewJSONClient(baseURL).
			WithRequestSender(requestSender)

		// Create a collection
		ctx := context.Background()
		err := client.CreateCollection(ctx, solr.NewCollectionParams().
			Name(collection).NumShards(1).ReplicationFactor(1))
		require.NoError(t, err, "creating a collection should not error")

		newTestRunner(client, collection)(t)

		// Delete the collection
		err = client.DeleteCollection(ctx, solr.NewCollectionParams().Name(collection))
		require.NoError(t, err, "deleting collection should not error")
	})
}

func newTestRunner(client solr.Client, collection string) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()

		// add new field type
		suggestText := solr.FieldType{
			Name:                 "suggest_text",
			Class:                "solr.TextField",
			PositionIncrementGap: "100",
			IndexAnalyzer: &solr.Analyzer{
				Tokenizer: &solr.Tokenizer{
					Class: "solr.WhitespaceTokenizerFactory",
				},
				Filters: []solr.Filter{
					{
						Class: "solr.LowerCaseFilterFactory",
					},
					{
						Class: "solr.ASCIIFoldingFilterFactory",
					},
					{
						Class:       "solr.EdgeNGramFilterFactory",
						MinGramSize: 1,
						MaxGramSize: 20,
					},
				},
			},
			QueryAnalyzer: &solr.Analyzer{
				Tokenizer: &solr.Tokenizer{
					Class: "solr.WhitespaceTokenizerFactory",
				},
				Filters: []solr.Filter{
					{
						Class: "solr.LowerCaseFilterFactory",
					},
					{
						Class: "solr.ASCIIFoldingFilterFactory",
					},
					{
						Class:    "solr.SynonymGraphFilterFactory",
						Synonyms: "synonyms.txt",
					},
				},
			},
		}
		err := client.AddFieldTypes(ctx, collection, suggestText)
		require.NoError(t, err, "adding field types should not error")

		// add fields
		fields := []solr.Field{
			{Name: "name", Type: "text_general"},
			{Name: "suggest", Type: "suggest_text"},
		}
		err = client.AddFields(ctx, collection, fields...)
		require.NoError(t, err, "adding fields should not error")

		// add copy fields
		copyFields := []solr.CopyField{
			{Source: "name", Dest: "suggest"},
			{Source: "name", Dest: "_text_"},
		}
		err = client.AddCopyFields(ctx, collection, copyFields...)
		require.NoError(t, err, "adding copy fields should not error")

		// Add suggester
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

		err = client.AddComponents(ctx, collection, suggestComponent, suggestHandler)
		require.NoError(t, err, "adding suggester components should not error")

		// Index
		docs := []solr.M{
			{"id": 1, "name": "Solr"},
			{"id": 2, "name": "Elastic"},
			{"id": 3, "name": "Blast"},
			{"id": 4, "name": "Bayard"},
		}
		buf := &bytes.Buffer{}
		err = json.NewEncoder(buf).Encode(docs)
		require.NoError(t, err, "encoding data should not error")

		_, err = client.Update(ctx, collection, solr.JSON, buf)
		require.NoError(t, err, "indexing data should not eror")

		err = client.Commit(ctx, collection)
		require.NoError(t, err, "commmit should not error")

		// Query
		qp := solr.NewStandardQueryParser().Query("*:*")
		query := solr.NewQuery(qp.BuildParser())
		queryResp, err := client.Query(ctx, collection, query)
		require.NoError(t, err, "query should not error")
		require.NotNil(t, queryResp, "query response should not be nil")
		assert.Len(t, queryResp.Response.Documents, 4, "query response is expected to have 4 documents")

		// Suggest
		queryStr := "solr"
		suggestParams := solr.NewSuggesterParams("suggest").Build().Query(queryStr)
		suggestResp, err := client.Suggest(ctx, collection, suggestParams)
		require.NoError(t, err, "suggest should not error")

		suggest := *suggestResp.Suggest
		termBody := suggest["default"][queryStr]
		assert.Len(t, termBody.Suggestions, 1, "expected to have one suggestion")
	}
}

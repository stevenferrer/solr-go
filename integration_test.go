// +build integration

package solr_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/sf9v/solr-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONClient(t *testing.T) {
	baseURL := "http://localhost:8983"
	collection := "searchengines"
	client := solr.NewJSONClient(baseURL)
	ctx := context.Background()

	// 1. add fields
	{
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
		require.NoError(t, err)

		fields := []solr.Field{
			{
				Name: "name",
				Type: "text_general",
			},
			{
				Name: "suggest",
				Type: "suggest_text",
			},
		}

		err = client.AddFields(ctx, collection, fields...)
		require.NoError(t, err)

		copyFields := []solr.CopyField{
			{
				Source: "name",
				Dest:   "suggest",
			},
			{
				Source: "name",
				Dest:   "_text_",
			},
		}

		err = client.AddCopyFields(ctx, collection, copyFields...)
		require.NoError(t, err)
	}
	// 2. add suggester component
	{
		suggestComponent := solr.NewComponent(solr.SearchComponent).
			Name("suggest").
			Class("solr.SuggestComponent").
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
			Name("/suggest").
			Class("solr.SearchHandler").
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
		require.NoError(t, err)
	}
	// 3. index
	{
		docs := []solr.M{
			{
				"id":   1,
				"name": "Solr",
			},
			{
				"id":   2,
				"name": "Elastic",
			},
			{
				"id":   3,
				"name": "Blast",
			},
			{
				"id":   4,
				"name": "Bayard",
			},
		}
		buf := &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(docs)
		require.NoError(t, err)

		_, err = client.Update(ctx, collection, solr.JSON, buf)
		require.NoError(t, err)

		err = client.Commit(ctx, collection)
		require.NoError(t, err)

		// send a request to build suggestions
		suggestParams := solr.NewSuggesterParams("suggest").Build()
		_, err = client.Suggest(ctx, collection, suggestParams)
		require.NoError(t, err)
	}
	// 4. query
	{
		queryParser := solr.NewStandardQueryParser().Query("*:*")
		query := solr.NewQuery().QueryParser(queryParser)
		queryResp, err := client.Query(ctx, collection, query)
		require.NoError(t, err)
		assert.NotNil(t, queryResp)
		assert.Len(t, queryResp.Response.Documents, 4)
	}

	// 5. query suggest endpoint
	{
		queryStr := "solr"
		suggestParams := solr.NewSuggesterParams("suggest").Query(queryStr)
		suggestResp, err := client.Suggest(ctx, collection, suggestParams)
		require.NoError(t, err)

		suggest := *suggestResp.Suggest
		termBody := suggest["default"][queryStr]
		assert.Len(t, termBody.Suggestions, 1)
	}
}

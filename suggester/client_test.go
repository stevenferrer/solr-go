package suggester_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stevenferrer/solr-go/config"
	"github.com/stevenferrer/solr-go/index"
	"github.com/stevenferrer/solr-go/schema"
	"github.com/stevenferrer/solr-go/suggester"
)

var (
	ctx             = context.Background()
	collection      = "gettingstarted"
	suggestEndpoint = "suggest"
	host            = "localhost"
	port            = 8983
	timeout         = time.Second * 60

	// only for covering
	_ = suggester.NewClient(host, port)
)

func TestClient(t *testing.T) {
	initSolr(t)

	t.Run("ok", func(t *testing.T) {
		rec, err := recorder.New("fixtures/suggest-ok")
		require.NoError(t, err)
		defer rec.Stop()

		client := suggester.NewCustomClient(host, port, suggestEndpoint, &http.Client{
			Timeout:   time.Second * 60,
			Transport: rec,
		})

		resp, err := client.Suggest(ctx, collection, suggester.Params{
			Query: "milana vino", Count: 10, Build: true,
		})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("handler not found", func(t *testing.T) {
			rec, err := recorder.New("fixtures/handler-not-found")
			require.NoError(t, err)
			defer rec.Stop()

			client := suggester.NewCustomClient(host, port, "not-exists",
				&http.Client{
					Timeout:   time.Second * 60,
					Transport: rec,
				},
			)

			_, err = client.Suggest(ctx, collection, suggester.Params{
				Query: "milana", Count: 10, Build: true,
			})
			assert.Error(t, err)
		})

		t.Run("parse url error", func(t *testing.T) {
			client := suggester.NewCustomClient("http\\\\\\::whatever:://\\::", 1234, "/not-exists",
				&http.Client{
					Timeout: time.Second * 60,
				},
			)

			_, err := client.Suggest(ctx, "gettingstarted///\\4343::343",
				suggester.Params{
					Query: "elec??&&",
					Count: 10,
					Build: true,
				})
			assert.Error(t, err)
		})

		t.Run("empty query", func(t *testing.T) {
			rec, err := recorder.New("fixtures/empty-query")
			require.NoError(t, err)
			defer rec.Stop()

			client := suggester.NewCustomClient(host, port, suggestEndpoint, &http.Client{
				Timeout:   time.Second * 60,
				Transport: rec,
			})

			_, err = client.Suggest(ctx, collection, suggester.Params{})
			assert.Error(t, err)
		})
	})
}

func initSolr(t *testing.T) {
	// initialize collection
	rec, err := recorder.New("fixtures/initialize")
	require.NoError(t, err)
	defer rec.Stop()

	schemaClient := schema.NewCustomClient(host, port, &http.Client{
		Timeout:   timeout,
		Transport: rec,
	})

	indexClient := index.NewCustomClient(host, port, &http.Client{
		Timeout:   timeout,
		Transport: rec,
	})

	configClient := config.NewCustomClient(host, port, &http.Client{
		Timeout:   timeout,
		Transport: rec,
	})

	// init schema
	err = schemaClient.AddFieldType(ctx, collection, schema.FieldType{
		Name:                 "text_suggest",
		Class:                "solr.TextField",
		PositionIncrementGap: "100",
		IndexAnalyzer: &schema.Analyzer{
			Tokenizer: &schema.Tokenizer{
				Class: "solr.StandardTokenizerFactory",
			},
			Filters: []schema.Filter{
				{
					Class: "solr.LowerCaseFilterFactory",
				},
				{
					Class:       "solr.EdgeNGramFilterFactory",
					MinGramSize: 1,
					MaxGramSize: 100,
				},
			},
		},
		QueryAnalyzer: &schema.Analyzer{
			Tokenizer: &schema.Tokenizer{
				Class: "solr.KeywordTokenizerFactory",
			},
			Filters: []schema.Filter{
				{
					Class: "solr.LowerCaseFilterFactory",
				},
			},
		},
	})
	require.NoError(t, err)

	err = schemaClient.AddField(ctx, collection, schema.Field{
		Name:    "name",
		Type:    "text_general",
		Indexed: true,
		Stored:  true,
	})
	require.NoError(t, err)

	// suggest field
	err = schemaClient.AddField(ctx, collection, schema.Field{
		Name:    "suggest",
		Type:    "text_suggest",
		Indexed: true,
		Stored:  true,
	})

	// add copy fields
	err = schemaClient.AddCopyField(ctx, collection, schema.CopyField{
		Source: "name",
		Dest:   "suggest",
	})
	require.NoError(t, err)

	// index some data
	var names = []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{
		{
			ID:   "1",
			Name: "Milana Vino",
		},
		{
			ID:   "2",
			Name: "Charly Jordan",
		},
		{
			ID:   "3",
			Name: "Daisy Keech",
		},
	}

	docs := index.NewDocs()
	for _, name := range names {
		docs.AddDoc(name)
	}

	err = indexClient.AddDocuments(ctx, collection, docs)
	require.NoError(t, err)

	// add suggester component
	addSuggestComponent := config.NewComponentCommand(
		config.AddSearchComponent,
		map[string]interface{}{
			"name":  suggestEndpoint,
			"class": "solr.SuggestComponent",
			"suggester": map[string]string{
				"name":                     "default",
				"lookupImpl":               "FuzzyLookupFactory",
				"dictionaryImpl":           "DocumentDictionaryFactory",
				"field":                    "suggest",
				"suggestAnalyzerFieldType": "text_suggest",
			},
		},
	)

	addSuggestHandler := config.NewComponentCommand(
		config.AddRequestHandler,
		map[string]interface{}{
			"name":    "/" + suggestEndpoint,
			"class":   "solr.SearchHandler",
			"startup": "lazy",
			"defaults": map[string]interface{}{
				"suggest":            true,
				"suggest.count":      10,
				"suggest.dictionary": "default",
			},
			"components": []string{"suggest"},
		},
	)

	err = configClient.SendCommands(ctx, collection,
		addSuggestComponent, addSuggestHandler)
	require.NoError(t, err)
}

package config_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	solrconfig "github.com/stevenferrer/solr-go/config"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	collection := "gettingstarted"
	host := "localhost"
	port := 8983
	timeout := time.Second * 6

	// only for covering
	_ = solrconfig.NewClient(host, port)

	t.Run("retrieve config", func(t *testing.T) {
		rec, err := recorder.New("fixtures/retrieve-config")
		require.NoError(t, err)
		defer rec.Stop()

		configClient := solrconfig.NewCustomClient(host, port, &http.Client{
			Timeout:   timeout,
			Transport: rec,
		})

		resp, err := configClient.GetConfig(ctx, collection)
		require.NoError(t, err)

		assert.NotNil(t, resp)
	})

	t.Run("send commands", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			rec, err := recorder.New("fixtures/send-commands-ok")
			require.NoError(t, err)
			defer rec.Stop()

			configClient := solrconfig.NewCustomClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			setUpdateHandlerAutoCommit := solrconfig.NewSetPropCommand(
				"updateHandler.autoCommit.maxTime", 15000)

			addSuggestComponent := solrconfig.NewComponentCommand(
				solrconfig.AddSearchComponent,
				map[string]interface{}{
					"name":  "suggest",
					"class": "solr.SuggestComponent",
					"suggester": map[string]string{
						"name":                     "mySuggester",
						"lookupImpl":               "FuzzyLookupFactory",
						"dictionaryImpl":           "DocumentDictionaryFactory",
						"field":                    "_text_",
						"suggestAnalyzerFieldType": "text_general",
					},
				},
			)

			addSuggestHandler := solrconfig.NewComponentCommand(
				solrconfig.AddRequestHandler,
				map[string]interface{}{
					"name":    "/suggest",
					"class":   "solr.SearchHandler",
					"startup": "lazy",
					"defaults": map[string]interface{}{
						"suggest":            true,
						"suggest.count":      10,
						"suggest.dictionary": "mySuggester",
					},
					"components": []string{"suggest"},
				},
			)

			err = configClient.SendCommands(ctx, collection,
				setUpdateHandlerAutoCommit,
				addSuggestComponent,
				addSuggestHandler,
			)
			assert.NoError(t, err)
		})

		t.Run("error", func(t *testing.T) {
			rec, err := recorder.New("fixtures/send-commands-error")
			require.NoError(t, err)
			defer rec.Stop()

			configClient := solrconfig.NewCustomClient(host, port, &http.Client{
				Timeout:   timeout,
				Transport: rec,
			})

			addSuggestComponent := solrconfig.NewComponentCommand(
				solrconfig.AddSearchComponent,
				map[string]interface{}{
					"name":  "suggest",
					"class": "solr.SuggestComponent",
					"suggester": map[string]string{
						"name":                     "mySuggester",
						"lookupImpl":               "FuzzyLookupFactory-BLAH-BLAH",
						"dictionaryImpl":           "DocumentDictionaryFactory-BLAH-BLAH",
						"field":                    "_text_",
						"suggestAnalyzerFieldType": "text_general",
					},
				},
			)

			err = configClient.SendCommands(ctx, collection, addSuggestComponent)
			assert.Error(t, err)
		})
	})

}

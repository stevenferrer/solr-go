package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	solrconfig "github.com/sf9v/solr-go/config"
)

func TestCommands(t *testing.T) {
	t.Run("common properties", func(t *testing.T) {
		t.Run("set property", func(t *testing.T) {
			setPropCommand := solrconfig.NewSetPropCommand("updateHandler.autoCommit.maxTime", 15000)

			got, err := setPropCommand.Command()
			require.NoError(t, err)

			expected := `"set-property": {"updateHandler.autoCommit.maxTime":15000}`
			assert.Equal(t, expected, got)
		})

		t.Run("unset property", func(t *testing.T) {
			unsetPropCommand := solrconfig.NewUnsetPropCommand("updateHandler.autoCommit.maxTime")

			got, err := unsetPropCommand.Command()
			require.NoError(t, err)

			expected := `"unset-property": "updateHandler.autoCommit.maxTime"`
			assert.Equal(t, expected, got)
		})
	})

	t.Run("component commands", func(t *testing.T) {
		// TODO: add table test for all command types
		addReqHandlerCommand := solrconfig.NewComponentCommand(
			solrconfig.AddRequestHandler,
			map[string]interface{}{
				"name":      "/mypath",
				"class":     "solr.DumpRequestHandler",
				"defaults":  map[string]interface{}{"x": "y", "a": "b", "rows": 10},
				"useParams": "x",
			},
		)

		got, err := addReqHandlerCommand.Command()
		require.NoError(t, err)

		expected := `"add-requesthandler": {"class":"solr.DumpRequestHandler","defaults":{"a":"b","rows":10,"x":"y"},"name":"/mypath","useParams":"x"}`
		assert.Equal(t, expected, got)
	})

	t.Run("add suggester command", func(t *testing.T) {
		dicts := []map[string]interface{}{
			{
				"name":                     "names",
				"lookupImpl":               "AnalyzingInfixLookupFactory",
				"dictionaryImpl":           "DocumentDictionaryFactory",
				"field":                    "_nameSuggest_",
				"weightField":              "price",
				"suggestAnalyzerFieldType": "suggest_text",
			},
			{
				"name":                     "tags",
				"lookupImpl":               "AnalyzingInfixLookupFactory",
				"dictionaryImpl":           "DocumentDictionaryFactory",
				"field":                    "_tagSuggest_",
				"weightField":              "price",
				"suggestAnalyzerFieldType": "suggest_text",
			},
		}
		addSuggester := solrconfig.NewAddSuggesterCommand("suggest", dicts...)

		got, err := addSuggester.Command()
		require.NoError(t, err)

		expected := `"add-searchcomponent":{"name":"suggest","class":"solr.SuggestComponent","suggester":{"dictionaryImpl":"DocumentDictionaryFactory","field":"_nameSuggest_","lookupImpl":"AnalyzingInfixLookupFactory","name":"names","suggestAnalyzerFieldType":"suggest_text","weightField":"price"}, "suggester":{"dictionaryImpl":"DocumentDictionaryFactory","field":"_tagSuggest_","lookupImpl":"AnalyzingInfixLookupFactory","name":"tags","suggestAnalyzerFieldType":"suggest_text","weightField":"price"}}`
		assert.Equal(t, expected, got)
	})
}

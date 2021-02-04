package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestBuildComponent(t *testing.T) {
	got := solr.NewComponent(solr.SearchComponent).
		Name("suggest").
		Class("solr.SearchComponent").
		Config(solr.M{
			"lookupImpl":               "AnalyzingInfixLookupFactory",
			"dictionaryImpl":           "DocumentDictionaryFactory",
			"field":                    "suggest",
			"suggestAnalyzerFieldType": "suggext_text",
		}).
		BuildComponent()

	expect := solr.M{
		"name":                     "suggest",
		"class":                    "solr.SearchComponent",
		"lookupImpl":               "AnalyzingInfixLookupFactory",
		"dictionaryImpl":           "DocumentDictionaryFactory",
		"field":                    "suggest",
		"suggestAnalyzerFieldType": "suggext_text",
	}

	assert.Equal(t, expect, got)
}

func TestComponentTypeStringer(t *testing.T) {
	var tests = []struct {
		componentType solr.ComponentType
		expected      string
	}{
		{
			solr.RequestHandler,
			"requesthandler",
		},
		{
			solr.SearchComponent,
			"searchcomponent",
		},
		{
			solr.InitParams,
			"initparams",
		},
		{
			solr.QueryResponseWriter,
			"queryresponsewriter",
		},
	}

	for _, test := range tests {
		got := test.componentType.String()
		assert.Equal(t, test.expected, got)
	}
}

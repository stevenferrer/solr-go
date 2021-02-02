package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestComponentTypeString(t *testing.T) {
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
		{
			solr.ComponentType(-1),
			"",
		},
	}

	for _, test := range tests {
		got := test.componentType.String()
		assert.Equal(t, test.expected, got)
	}
}

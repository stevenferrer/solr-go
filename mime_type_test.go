package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestContentTypeStringer(t *testing.T) {
	var tests = []struct {
		mimeType solr.MimeType
		expected string
	}{
		{
			solr.JSON,
			"application/json",
		},
		{
			solr.XML,
			"application/xml",
		},
		{
			solr.CSV,
			"text/csv",
		},
	}

	for _, test := range tests {
		got := test.mimeType.String()
		assert.Equal(t, test.expected, got)
	}
}

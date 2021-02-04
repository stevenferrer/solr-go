package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestContentTypeStringer(t *testing.T) {
	var tests = []struct {
		contentType solr.ContentType
		expected    string
	}{
		{
			solr.JSON,
			"application/json",
		},
		{
			solr.XML,
			"application/xml",
		},
	}

	for _, test := range tests {
		got := test.contentType.String()
		assert.Equal(t, test.expected, got)
	}
}

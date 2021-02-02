package solr_test

import (
	"testing"

	"github.com/sf9v/solr-go"
	"github.com/stretchr/testify/assert"
)

func TestSuggesterParams(t *testing.T) {
	got := solr.NewSuggesterParams("/suggest").
		Query("electronic dev").
		Dictionaries("default", "custom").
		Count(10).Cfq("memory").Build().
		Reload().BuildAll().ReloadAll().
		BuildParams()

	expect := `suggest=true&suggest.build=true&suggest.buildAll=true&suggest.cfg=memory&suggest.count=10&suggest.dictionary=default&suggest.dictionary=custom&suggest.q=electronic+dev&suggest.reload=true&suggest.reloadAll=true`
	assert.Equal(t, expect, got)
}

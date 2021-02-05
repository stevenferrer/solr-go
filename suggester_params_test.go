package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestBuildSuggesterParams(t *testing.T) {
	got := solr.NewSuggesterParams("suggest").
		Query("electronic dev").
		Dictionaries("default", "custom").
		Count(10).Cfq("memory").Build().
		Reload().BuildAll().ReloadAll().
		BuildParams()

	expect := `suggest=true&suggest.build=true&suggest.buildAll=true&suggest.cfg=memory&suggest.count=10&suggest.dictionary=default&suggest.dictionary=custom&suggest.q=electronic+dev&suggest.reload=true&suggest.reloadAll=true`
	assert.Equal(t, expect, got)

	got = solr.NewSuggesterParams("suggest").Build().BuildParams()
	expect = `suggest=true&suggest.build=true`
	assert.Equal(t, expect, got)
}

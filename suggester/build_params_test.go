package suggester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildURLParams(t *testing.T) {
	got := buildURLParams(Params{
		Query:        "electronic dev",
		Dictionaries: []string{"default", "customDict"},
		Count:        10,
		Cfq:          "memory",
		Build:        true,
		Reload:       true,
		BuildAll:     true,
		ReloadAll:    true,
	})

	expect := `suggest=true&suggest.build=true&suggest.buildAll=true&suggest.cfg=memory&suggest.count=10&suggest.dictionary=default&suggest.dictionary=customDict&suggest.q=electronic+dev&suggest.reload=true&suggest.reloadAll=true`
	assert.Equal(t, expect, got)
}

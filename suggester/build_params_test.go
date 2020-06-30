package suggester

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildParams(t *testing.T) {
	params := buildParams(Params{
		Query:        "elec",
		Dictionaries: []string{"mySuggester"},
		Count:        10,
		Cfq:          "memory",
		Build:        true,
		Reload:       true,
		BuildAll:     true,
		ReloadAll:    true,
	})

	assert.Len(t, params, 9)
}

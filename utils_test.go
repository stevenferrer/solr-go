package solr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeQueryChars(t *testing.T) {
	a := assert.New(t)
	queryFromDocs := "(1+1):2"
	expected := `\(1\+1\)\:2`
	a.Equal(expected, EscapeQueryChars(queryFromDocs))

	queryAllSpecialChars := `\+-!():^[]"{}~*?|&;/`
	expected2 := `\\\+\-\!\(\)\:\^\[\]\"\{\}\~\*\?\|\&\;\/`
	a.Equal(expected2, EscapeQueryChars(queryAllSpecialChars))

	queryWhitespace := `	solr rocks`
	expected3 := `\	solr\ rocks`
	a.Equal(expected3, EscapeQueryChars(queryWhitespace))
}

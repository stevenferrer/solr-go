package solr

import (
	"strings"
	"unicode"
)

// EscapeQueryChars escapes special characters in the given Solr query string that would normally be treated as part of
// Solr's query syntax. For a full list of special characters, see the Solr documentation here:
// https://solr.apache.org/guide/solr/9_7/query-guide/standard-query-parser.html#escaping-special-characters
//
// Returns the query string with special characters escaped
func EscapeQueryChars(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '\\', '+', '-', '!', '(', ')', ':', '^', '[', ']', '"', '{', '}', '~', '*', '?', '|', '&', ';', '/':
			sb.WriteRune('\\')
		}
		if unicode.IsSpace(c) {
			sb.WriteRune('\\')
		}
		sb.WriteRune(c)
	}
	return sb.String()
}

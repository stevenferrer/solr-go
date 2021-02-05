// Package solr is a solr client for Go
package solr

// ContentType is a content-type
type ContentType int

// List of content-types
const (
	JSON ContentType = iota
	XML
)

// String implements Stringer
func (ct ContentType) String() string {
	return [...]string{
		"application/json",
		"application/xml",
	}[ct]
}

// M is an alias for map of interface
type M map[string]interface{}

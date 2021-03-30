package solr

// MimeType is a mime-type
type MimeType int

// List of mime-types
const (
	JSON MimeType = iota
	XML
	CSV
)

// String implements Stringer
func (mt MimeType) String() string {
	return [...]string{
		"application/json",
		"application/xml",
		"text/csv",
	}[mt]
}

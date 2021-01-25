package solr

// Facets is a collection of facets
type Facets []Faceter

// Faceter is an abstraction of a facet
// e.g. terms, query, stats, heatmap etc.
type Faceter interface {
	BuildFacet()
}

// TermsFacet is a terms facet
type TermsFacet struct {
	field       string
	offset      int
	limit       int
	sort        string
	overRequest int
	refine      bool
	overRefine  int
	minCount    int
	missing     bool
	numBuckets  bool
	allBuckets  bool
	prefix      string
	method      string
	prelimSort  string
	facet       M
	// Nested facets
	// https://lucene.apache.org/solr/guide/8_7/json-facet-api.html#nested-facets
	facets Facets
}

// QueryFacet is query facet
type QueryFacet struct {
	q string
}

// RangeFacet is a range facet
type RangeFacet struct {
	field   string
	start   int
	end     int
	gap     int
	hardend bool
	other   string
	include string
	ranges  string
	facet   M
	// Nested facets
	facets Facets
}

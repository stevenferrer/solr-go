package solr

// Facets is a collection of facets
type Facets []Faceter

// Faceter is an abstraction of a facet
// e.g. terms, query, stats, heatmap etc.
type Faceter interface {
	// BuildFacet builds the facet
	BuildFacet() M
	// Name gives the name of the facet
	Name() string
}

// TermsFacet is a terms facet
type TermsFacet struct {
	// name is the facet name
	name string

	// terms facet params
	field  string
	offset int
	limit  int
	sort   string

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

	facet M
}

var _ Faceter = (*TermsFacet)(nil)

// NewTermsFacet returns a new terms facet
func NewTermsFacet(name string) *TermsFacet {
	return &TermsFacet{name: name, facet: M{}}
}

// BuildFacet builds the facet
func (f *TermsFacet) BuildFacet() M {
	m := M{"type": "terms"}

	if f.field != "" {
		m["field"] = f.field
	}

	if f.offset > 0 {
		m["offset"] = f.offset
	}

	if f.limit > 0 {
		m["limit"] = f.limit
	}

	if f.sort != "" {
		m["sort"] = f.sort
	}

	if len(f.facet) > 0 {
		m["facet"] = f.facet
	}

	return m
}

// WithField sets the field param
func (f *TermsFacet) WithField(field string) *TermsFacet {
	f.field = field
	return f
}

// WithOffset sets the offset param
func (f *TermsFacet) WithOffset(offset int) *TermsFacet {
	f.offset = offset
	return f
}

// WithLimit sets the limit param
func (f *TermsFacet) WithLimit(limit int) *TermsFacet {
	f.limit = limit
	return f
}

// WithSort sets the sort param
func (f *TermsFacet) WithSort(sort string) *TermsFacet {
	f.sort = sort
	return f
}

// AddNestedFacet adds a nested facet
func (f *TermsFacet) AddNestedFacet(facet Faceter) *TermsFacet {
	f.facet[facet.Name()] = facet.BuildFacet()
	return f
}

// AddToFacet adds a key-value pair to the facet map
func (f *TermsFacet) AddToFacet(key string, value interface{}) *TermsFacet {
	f.facet[key] = value
	return f
}

// Name returns the name of the facet
func (f *TermsFacet) Name() string {
	return f.name
}

// QueryFacet is query facet
type QueryFacet struct {
	name  string
	q     string
	facet M
}

var _ Faceter = (*QueryFacet)(nil)

// NewQueryFacet retuns a new query facet
func NewQueryFacet(name string) *QueryFacet {
	return &QueryFacet{name: name, facet: M{}}
}

// BuildFacet builds the facet
func (f *QueryFacet) BuildFacet() M {
	m := M{"type": "query"}
	if f.q != "" {
		m["q"] = f.q
	}

	if len(f.facet) > 0 {
		m["facet"] = f.facet
	}

	return m
}

// Name returns the facet name
func (f *QueryFacet) Name() string {
	return f.name
}

// WithQuery sets the q param
func (f *QueryFacet) WithQuery(q string) *QueryFacet {
	f.q = q
	return f
}

// AddNestedFacet adds a nested facet
func (f *QueryFacet) AddNestedFacet(facet Faceter) *QueryFacet {
	f.facet[facet.Name()] = facet.BuildFacet()
	return f
}

// AddToFacet adds a key-value pair to the facet map
func (f *QueryFacet) AddToFacet(key string, value interface{}) *QueryFacet {
	f.facet[key] = value
	return f
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

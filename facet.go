package solr

// Faceter is an abstraction of a facet
// e.g. terms, query, stats, range, heatmap etc.
type Faceter interface {
	// BuildFacet builds the facet
	BuildFacet() M
	// Name returns the name of the facet
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

	// overRequest int
	// refine      bool
	// overRefine  int
	minCount int
	// missing     bool
	// numBuckets  bool
	// allBuckets  bool
	// prefix      string
	// method      string
	// prelimSort  string

	facet M

	domain M
}

var _ Faceter = (*TermsFacet)(nil)

// NewTermsFacet returns a new TermsFacet
func NewTermsFacet(name string) *TermsFacet {
	return &TermsFacet{name: name, facet: M{}, domain: M{}}
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

	if f.minCount > 0 {
		m["mincount"] = f.minCount
	}

	if len(f.facet) > 0 {
		m["facet"] = f.facet
	}

	if len(f.domain) > 0 {
		m["domain"] = f.domain
	}

	return m
}

// Field sets the field param
func (f *TermsFacet) Field(field string) *TermsFacet {
	f.field = field
	return f
}

// Offset sets the offset param
func (f *TermsFacet) Offset(offset int) *TermsFacet {
	f.offset = offset
	return f
}

// Limit sets the limit param
func (f *TermsFacet) Limit(limit int) *TermsFacet {
	f.limit = limit
	return f
}

// Sort sets the sort param
func (f *TermsFacet) Sort(sort string) *TermsFacet {
	f.sort = sort
	return f
}

// AddFacets adds nested facets
func (f *TermsFacet) AddFacets(facets ...Faceter) *TermsFacet {
	for _, facet := range facets {
		f.facet[facet.Name()] = facet.BuildFacet()
	}

	return f
}

// AddToFacet adds a key-value pair to the facet map
func (f *TermsFacet) AddToFacet(key string, value interface{}) *TermsFacet {
	f.facet[key] = value
	return f
}

// AddToDomain adds a key-value pair to the domain map
func (f *TermsFacet) AddToDomain(key string, value interface{}) *TermsFacet {
	f.domain[key] = value
	return f
}

// Name returns the name of the facet
func (f *TermsFacet) Name() string {
	return f.name
}

// MinCount sets the mincount param
func (f *TermsFacet) MinCount(minCount int) *TermsFacet {
	f.minCount = minCount
	return f
}

// QueryFacet is query facet
type QueryFacet struct {
	name  string
	q     string
	facet M
}

var _ Faceter = (*QueryFacet)(nil)

// NewQueryFacet retunrs a new QueryFacet
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

// Query sets the q param
func (f *QueryFacet) Query(q string) *QueryFacet {
	f.q = q
	return f
}

// AddFacet adds a nested facet
func (f *QueryFacet) AddFacet(facet Faceter) *QueryFacet {
	f.facet[facet.Name()] = facet.BuildFacet()
	return f
}

// AddToFacet adds a key-value pair to the facet map
func (f *QueryFacet) AddToFacet(key string, value interface{}) *QueryFacet {
	f.facet[key] = value
	return f
}

// // RangeFacet is a range facet
// type RangeFacet struct {
// 	field   string
// 	start   int
// 	end     int
// 	gap     int
// 	hardened bool
// 	other   string
// 	include string
// 	ranges  string
// 	facet   M
// }

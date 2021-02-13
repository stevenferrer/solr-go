package solr

import (
	"net/url"
	"strconv"
)

// SuggestParams is the suggester API param builder
type SuggestParams struct {
	endpoint string

	// dicts is the name of the dictionary
	// component configured in the search component
	dicts []string

	// query is the query to use for suggestion lookups
	query,

	// cfq is a context filter query used to filter suggestsion
	// based on the context field, if supported by the suggester
	cfq string

	// count is the number of suggestions to return
	count int

	// build If true, it will build the suggester index
	build,

	// Reload If true, it will reload the suggester index.
	reload,

	// BuildAll If true, it will build all suggester indexes.
	buildAll,

	// ReloadAll If true, it will reload all suggester indexes.
	reloadAll bool
}

// NewSuggesterParams returns a new SuggesterParams
func NewSuggesterParams(endpoint string) *SuggestParams {
	return &SuggestParams{endpoint: endpoint}
}

// BuildParams builds the suggester params
func (sp *SuggestParams) BuildParams() string {
	params := url.Values{}

	params.Add("suggest", "true")

	for _, dict := range sp.dicts {
		params.Add("suggest.dictionary", dict)
	}

	if sp.query != "" {
		params.Add("suggest.q", sp.query)
	}

	if sp.count > 0 {
		params.Add("suggest.count", strconv.Itoa(sp.count))
	}

	if sp.cfq != "" {
		params.Add("suggest.cfg", sp.cfq)
	}

	if sp.build {
		params.Add("suggest.build", "true")
	}

	if sp.reload {
		params.Add("suggest.reload", "true")
	}

	if sp.buildAll {
		params.Add("suggest.buildAll", "true")
	}

	if sp.reloadAll {
		params.Add("suggest.reloadAll", "true")
	}

	return params.Encode()
}

// Dictionaries sets the dictionaries param
func (sp *SuggestParams) Dictionaries(dicts ...string) *SuggestParams {
	sp.dicts = dicts
	return sp
}

// Query sets the query param
func (sp *SuggestParams) Query(query string) *SuggestParams {
	sp.query = query
	return sp
}

// Cfq sets the context filter query param
func (sp *SuggestParams) Cfq(cfq string) *SuggestParams {
	sp.cfq = cfq
	return sp
}

// Count sets the count param
func (sp *SuggestParams) Count(count int) *SuggestParams {
	sp.count = count
	return sp
}

// Build sets the build param to true
func (sp *SuggestParams) Build() *SuggestParams {
	sp.build = true
	return sp
}

// Reload sets the reload param to true
func (sp *SuggestParams) Reload() *SuggestParams {
	sp.reload = true
	return sp
}

// BuildAll sets the build-all param to true
func (sp *SuggestParams) BuildAll() *SuggestParams {
	sp.buildAll = true
	return sp
}

// ReloadAll sets the reload-all param to true
func (sp *SuggestParams) ReloadAll() *SuggestParams {
	sp.reloadAll = true
	return sp
}

package suggester

// Request is a suggest request
type Request struct {
	// Collection is the name of the collection
	Collection string
	// Params is the suggester params
	Params Params
}

// Params is the suggester parameters
type Params struct {
	// Dictionaries is the name of the dictionary
	// component configured in the search component
	Dictionaries []string

	// Query is the query to use for suggestion lookups
	Query string

	// Count is the number of suggestions to return
	Count int

	// ContextFilterQuery  is a context filter query used to filter
	// suggestsion based on the context field, if supported by the suggester
	Cfq string

	// Build If true, it will build the suggester index
	Build,

	// Reload If true, it will reload the suggester index.
	Reload,

	// BuildAll If true, it will build all suggester indexes.
	BuildAll,

	// ReloadAll If true, it will reload all suggester indexes.
	ReloadAll bool
}

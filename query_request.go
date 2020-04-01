package helios

// JSONer ...
type JSONer interface {
	JSON() ([]byte, error)
}

// M is an alias to map of interfaces
type M map[string]interface{}

// SimpleQueryRequest is used for simple query
// request where the query is only a string
// {
//   "query" : "name:iPod"
// }
type SimpleQueryRequest struct {
	Query   string   `json:"query"`
	Limit   uint     `json:"limit,omitempty"`
	Filter  []string `json:"filter,omitempty"`
	Facet   M        `json:"facet,omitempty"`
	Queries M        `json:"queries,omitempty"`
}

// ExpandedQueryRequest is used for more advanced query
// request where the query is a fully expanded JSON object
//
// Lucene query:
// {
//   "query": {
//     "lucene": {
//       "df": "name",
//       "query": "iPod"
//     }
//   }
// }
//
// Boost query:
// {
//   "query": {
//     "boost": {
//       "query": {
//         "lucene": {
//           "df": "name",
//           "query": "iPod"
//         }
//       },
//       "b": "log(popularity)"
//     }
//   }
// }
type ExpandedQueryRequest struct {
	Query   M        `json:"query"`
	Limit   uint     `json:"limit,omitempty"`
	Filter  []string `json:"filter,omitempty"`
	Facet   M        `json:"facet,omitempty"`
	Queries M        `json:"queries,omitempty"`
}

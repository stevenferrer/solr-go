package solr

// QueryResponse is a query response
type QueryResponse struct {
	BaseResponse
	Response QueryResponseBody `json:"response,omitempty"`
	Facets   M                 `json:"facets,omitempty"`
}

// QueryResponseBody is the query response body
type QueryResponseBody struct {
	NumFound int     `json:"numFound,omitempty"`
	Start    int     `json:"start,omitempty"`
	MaxScore float64 `json:"maxScore,omitempty"`
	Docs     []M     `json:"docs,omitempty"`
}

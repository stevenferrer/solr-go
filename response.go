package solr

// BaseResponse is the base response
type BaseResponse struct {
	Header *ResponseHeader `json:"responseHeader"`
	Error  *Error          `json:"error,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	ZKConnected bool `json:"zkConnected"`
	Status      int  `json:"status"`
	QTime       int  `json:"QTime"`
}

// Error is a response error
type Error struct {
	Code     int      `json:"code"`
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

// UpdateResponse is an update response
type UpdateResponse struct {
	*BaseResponse
}

// QueryResponse is a query response
type QueryResponse struct {
	*BaseResponse
	Response QueryResponseBody `json:"response,omitempty"`
	Facets   M                 `json:"facets,omitempty"`
}

// QueryResponseBody is the query response body
type QueryResponseBody struct {
	NumFound  int     `json:"numFound,omitempty"`
	Start     int     `json:"start,omitempty"`
	MaxScore  float64 `json:"maxScore,omitempty"`
	Documents []M     `json:"docs,omitempty"`
}

// SuggestResponse is the suggester response
type SuggestResponse struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Command        string         `json:"command,omitempty"`
	Suggest        *SuggestBody   `json:"suggest,omitempty"`
	Error          *Error         `json:"error,omitempty"`
}

// SuggestBody is the suggest body
type SuggestBody map[string]SuggestTerm

// SuggestTerm is suggest term
type SuggestTerm map[string]SuggestTermBody

// SuggestTermBody is the suggest term body
type SuggestTermBody struct {
	NumFound    int          `json:"numFound,omitempty"`
	Suggestions []Suggestion `json:"suggestions,omitempty"`
}

// Suggestion is a term suggestion
type Suggestion struct {
	Term    string `json:"term"`
	Weight  int    `json:"weight"`
	Payload string `json:"payload,omitempty"`
}

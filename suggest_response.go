package solr

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

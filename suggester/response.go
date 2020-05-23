package suggester

// Response is the suggester response
type Response struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Command        string         `json:"command,omitempty"`
	Suggest        *SuggestBody   `json:"suggest,omitempty"`
	Error          *Error         `json:"error,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

// SuggestBody is the suggest body
// Key: Suggester dictionary
// Value: SuggestTerm
type SuggestBody map[string]SuggestTerm

// SuggestTerm is suggest term
// Key: Suggest term query
// Value: SuggestTermBody
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

// Error is a response error
type Error struct {
	Code     int      `json:"code"`
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

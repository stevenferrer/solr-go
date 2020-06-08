package query

import (
	"github.com/stevenferrer/solr-go/types"
)

// Response is a query response
type Response struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       ResponseBody   `json:"response,omitempty"`
	Facets         types.M        `json:"facets,omitempty"`
	Error          *Error         `json:"error,omitempty"`
}

// ResponseBody is the response body
type ResponseBody struct {
	NumFound int       `json:"numFound,omitempty"`
	Start    int       `json:"start,omitempty"`
	MaxScore float64   `json:"maxScore,omitempty"`
	Docs     []types.M `json:"docs,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
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

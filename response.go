package solr

import "time"

// BaseResponse is the base response
type BaseResponse struct {
	Header *ResponseHeader `json:"responseHeader"`
	Error  *ResponseError  `json:"error,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	ZKConnected bool `json:"zkConnected"`
	Status      int  `json:"status"`
	QTime       int  `json:"QTime"`
}

// ResponseError is a response error
type ResponseError struct {
	Code     int      `json:"code"`
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
}

func (e ResponseError) Error() string {
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
	*BaseResponse
	Command string       `json:"command,omitempty"`
	Suggest *SuggestBody `json:"suggest,omitempty"`
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

type CoreStatusResponse struct {
	*BaseResponse
	InitFailures M                      `json:"initFailures"`
	Status       map[string]*CoreStatus `json:"status"`
}

type CoreStatus struct {
	Config      string     `json:"config"`
	DataDir     string     `json:"dataDir"`
	Index       *Index     `json:"index"`
	InstanceDir string     `json:"instanceDir"`
	Name        string     `json:"name"`
	Schema      string     `json:"schema"`
	StartTime   *time.Time `json:"startTime"`
	Uptime      int        `json:"uptime"`
}

type Index struct {
	Current                 bool   `json:"current"`
	DeletedDocs             int    `json:"deletedDocs"`
	Directory               string `json:"directory"`
	HasDeletions            bool   `json:"hasDeletions"`
	IndexHeapUsageBytes     int    `json:"indexHeapUsageBytes"`
	MaxDoc                  int    `json:"maxDoc"`
	NumDocs                 int    `json:"numDocs"`
	SegmentCount            int    `json:"segmentCount"`
	SegmentFile             string `json:"segmentFile"`
	SegmentsFileSizeInBytes int    `json:"segmentsFileSizeInBytes"`
	Size                    string `json:"size"`
	UserData                M      `json:"userData"`
	Version                 int
}

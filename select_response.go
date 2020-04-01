package helios

// SelectResponse is the select response
type SelectResponse struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Response       *Response      `json:"response,omitempty"`
	Facets         Facets         `json:"facets,omitempty"`
}
type ResponseHeader struct {
	ZkConnected bool    `json:"zkConnected"`
	Status      int     `json:"status"`
	QTime       int     `json:"QTime"`
	Params      *Params `json:"params,omitempty"`
}

type Params struct {
	JSON       string `json:"json,omitempty"`
	Q          string `json:"q,omitempty"`
	FacetField string `json:"facet.field,omitempty"`
	Facet      string `json:"facet,omitempty"`
	Underscore string `json:"_,omitempty"`
}

type Response struct {
	NumFound int     `json:"numFound"`
	Start    int     `json:"start,omitempty"`
	MaxScore float64 `json:"maxScore,omitempty"`
	Docs     []M     `json:"docs"`
}
type Buckets struct {
	Val   string `json:"val"`
	Count int    `json:"count"`
}
type Categories struct {
	Buckets []Buckets `json:"buckets"`
}
type Facets struct {
	Count      int        `json:"count"`
	Categories Categories `json:"categories"`
}

package solr

// BaseResponse is the base response
type BaseResponse struct {
	Header *ResponseHeader `json:"responseHeader"`
	Error  *Error          `json:"error,omitempty"`
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

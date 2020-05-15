package index

// Response is an update response
type Response struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Error          *Error         `json:"error,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

// Error is the response error detail
type Error struct {
	Code     int      `json:"code"`
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

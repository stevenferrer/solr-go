package helios

type ResponseHeader struct {
	ZKConnected bool    `json:"zkConnected,omitempty"`
	Status      int     `json:"status"`
	QTime       int     `json:"QTime"`
	Params      *Params `json:"params,omitempty"`
}

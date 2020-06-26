package config

// Response is a config response
type Response struct {
	ResponseHeader *ResponseHeader `json:"responseHeader,omitempty"`
	Config         *Config         `json:"config,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

// Config response body
type Config struct {
	LuceneMatchVersion          string                 `json:"luceneMatcheVersion"`
	UpdateHandler               map[string]interface{} `json:"updateHandler,omitempty"`
	Query                       map[string]interface{} `json:"query,omitempty"`
	RequestHandler              map[string]interface{} `json:"requestHandler,omitempty"`
	QueryResponseWriter         map[string]interface{} `json:"queryResponseWriter,omitempty"`
	SearchComponent             map[string]interface{} `json:"searchComponent,omitempty"`
	UpdateProcessor             map[string]interface{} `json:"updateProcessor,omitempty"`
	InitParams                  []interface{}          `json:"initParams,omitempty"`
	Listener                    []interface{}          `json:"listener,omitempty"`
	DirectoryFactory            map[string]interface{} `json:"directoryFactory,omitempty"`
	CodeFactory                 map[string]interface{} `json:"codeFactory,omitempty"`
	UpdateRequestProcessorChain []interface{}          `json:"updateRequestProcessorChain,omitempty"`
	UpdateHandlerUpdateLog      map[string]interface{} `json:"updateHandlerupdateLog,omitempty"`
	RequestDispatcher           map[string]interface{} `json:"requestDispatcher,omitempty"`
	IndexConfig                 map[string]interface{} `json:"indexCofig,omitempty"`
	PeerSync                    map[string]interface{} `json:"peerSync,omitempty"`
}

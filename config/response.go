package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Response is a config response
type Response struct {
	ResponseHeader *ResponseHeader `json:"responseHeader,omitempty"`
	Config         *Config         `json:"config,omitempty"`
	Error          *Error          `json:"error,omitempty"`
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

// Error is the response error detail
type Error struct {
	Code    int `json:"code"`
	Details []struct {
		command

		ErrorMessages []string `json:"errorMessages"`
	} `json:"details"`
	Metadata []string `json:"metadata"`
	Msg      string   `json:"msg"`
}

type command struct {
	AddRequestHandler         interface{} `json:"add-requesthandler"`
	UpdateRequestHandler      interface{} `json:"update-requesthandler"`
	DeleteRequestHandler      interface{} `json:"delete-requesthandler"`
	AddSearchComponent        interface{} `json:"add-searchcomponent"`
	UpdateSearchComponent     interface{} `json:"update-searchcomponent"`
	DeleteSearchComponent     interface{} `json:"delete-searchcomponent"`
	AddInitParams             interface{} `json:"add-initparams"`
	UpdateInitParams          interface{} `json:"update-initparams"`
	DeleteInitParams          interface{} `json:"delete-initparams"`
	AddQueryResponseWriter    interface{} `json:"add-queryresponsewriter"`
	UpdateQueryResponseWriter interface{} `json:"update-queryresponsewriter"`
	DeleteQueryResponseWriter interface{} `json:"delete-queryresponsewriter"`
}

func (c command) getCmd() string {
	cv := reflect.ValueOf(c)
	for i := 0; i < cv.NumField(); i++ {
		cf := cv.Field(i)
		if cf.IsZero() {
			continue
		}

		tagVal := cv.Type().Field(i).Tag.Get("json")
		args := strings.Split(tagVal, ",")
		return args[0]
	}

	return "unknown"
}

func (e Error) Error() string {
	errMsgs := []string{}

	for _, det := range e.Details {

		errs := []string{}
		for _, errMsg := range det.ErrorMessages {
			errs = append(errs, errMsg)
		}

		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s",
			det.getCmd(), strings.Join(errs, ", ")))

	}

	return fmt.Sprintf("%s: %s", e.Msg, strings.Join(errMsgs, ", "))
}

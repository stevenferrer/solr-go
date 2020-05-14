package schema

import (
	"fmt"
	"reflect"
	"strings"
)

// Response is a schema API response
type Response struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Error          *Error         `json:"error,omitempty"`
}

// ResponseHeader is a response header
type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

type command struct {
	AddField            interface{} `json:"add-field"`
	DeleteField         interface{} `json:"delete-field"`
	ReplaceField        interface{} `json:"replace-field"`
	AddDynamicField     interface{} `json:"add-dynamic-field"`
	DeleteDynamicField  interface{} `json:"delete-dynamic-field"`
	ReplaceDynamicField interface{} `json:"replace-dynamic-field"`
	AddCopyField        interface{} `json:"add-copy-field"`
	DeleteCopyField     interface{} `json:"delete-copy-field"`
	AddFieldType        interface{} `json:"add-field-type"`
	ReplaceFieldType    interface{} `json:"replace-field-type"`
	DeleteFieldType     interface{} `json:"delete-field-type"`
}

func (c command) GetCommand() string {
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

	return ""
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

func (e Error) Error() string {
	errMsgs := []string{}

	for _, det := range e.Details {

		errs := []string{}
		for _, errMsg := range det.ErrorMessages {
			errs = append(errs, errMsg)
		}

		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s",
			det.GetCommand(), strings.Join(errs, ", ")))

	}

	return fmt.Sprintf("%s: %s", e.Msg, strings.Join(errMsgs, ", "))
}

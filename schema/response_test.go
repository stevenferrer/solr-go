package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name    string
		respErr []byte
		want    string
	}{
		{
			name: "add field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "add-copy-field":{},
              "errorMessages":["'source' is a required field",
                "'dest' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: add-copy-field: 'source' is a required field, 'dest' is a required field",
		},
		{
			name: "add dynamic field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "add-dynamic-field":{},
              "errorMessages":["'name' is a required field",
                "'type' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: add-dynamic-field: 'name' is a required field, 'type' is a required field",
		},
		{
			name: "add field type error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "add-field-type":{},
              "errorMessages":["'name' is a required field",
                "'class' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: add-field-type: 'name' is a required field, 'class' is a required field",
		},
		{
			name: "add field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "add-field":{},
              "errorMessages":["'name' is a required field",
                "'type' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: add-field: 'name' is a required field, 'type' is a required field",
		},
		{
			name: "delete copy field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "delete-copy-field":{},
              "errorMessages":["'source' is a required field",
                "'dest' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: delete-copy-field: 'source' is a required field, 'dest' is a required field",
		},
		{
			name: "delete dynamic field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "delete-dynamic-field":{},
              "errorMessages":["'name' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: delete-dynamic-field: 'name' is a required field",
		},
		{
			name: "delete field type error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "delete-field-type":{},
              "errorMessages":["'name' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: delete-field-type: 'name' is a required field",
		},
		{
			name: "delete field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "delete-field":{},
              "errorMessages":["'name' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: delete-field: 'name' is a required field",
		},
		{
			name: "replace dynamic field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "replace-dynamic-field":{},
              "errorMessages":["'name' is a required field",
                "'type' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: replace-dynamic-field: 'name' is a required field, 'type' is a required field",
		},
		{
			name: "replace field type error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "replace-field-type":{},
              "errorMessages":["'name' is a required field",
                "'class' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: replace-field-type: 'name' is a required field, 'class' is a required field",
		},
		{
			name: "replace field error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "replace-field":{},
              "errorMessages":["'name' is a required field",
                "'type' is a required field"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: replace-field: 'name' is a required field, 'type' is a required field",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			var e Error
			err := json.Unmarshal(tt.respErr, &e)
			require.NoError(t, err)

			got := e.Error()
			a.Equal(tt.want, got)
		})
	}
}

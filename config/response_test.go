package config

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
			name: "add search component error",
			respErr: []byte(`{
          "metadata":[
            "error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject",
            "root-error-class","org.apache.solr.api.ApiBag$ExceptionWithErrObject"],
          "details":[{
              "add-searchcomponent":{
                "class":"solr.SuggestComponent",
                "name":"suggest",
                "suggester":{
                  "dictionaryImpl":"DocumentDictionaryFactory-BLAH-BLAH",
                  "field":"_text_",
                  "lookupImpl":"FuzzyLookupFactory-BLAH-BLAH",
                  "name":"mySuggester",
                  "suggestAnalyzerFieldType":"text_general"}},
              "errorMessages":[" 'suggest' already exists . Do an 'update-searchcomponent' , if you want to change it "]}],
          "msg":"error processing commands",
          "code":400}`),
			want: "error processing commands: add-searchcomponent:  'suggest' already exists . Do an 'update-searchcomponent' , if you want to change it ",
		},
		{
			name: "unknown command",
			respErr: []byte(`{
          "details":[{
              "add-queryparser":{},
              "errorMessages":["an error"]}],
          "msg":"error processing commands",
		  "code":400}`),
			want: "error processing commands: unknown: an error",
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

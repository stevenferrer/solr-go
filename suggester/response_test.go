package suggester_test

import (
	"encoding/json"
	"testing"

	"github.com/sf9v/solr-go/suggester"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalResponse(t *testing.T) {
	jsonStr := `{
		"responseHeader": {
			"status": 0,
			"QTime": 3
		},
		"command": "build",
		"suggest": {
			"mySuggester": {
			"elec": {
				"numFound": 1,
				"suggestions": [
				{
					"term": "electronics and computer1",
					"weight": 100,
					"payload": ""
				}
				]
			}
			},
			"altSuggester": {
			"elec": {
				"numFound": 1,
				"suggestions": [
				{
					"term": "electronics and computer1",
					"weight": 10,
					"payload": ""
				}
				]
			}
			}
		}
	}`

	var response suggester.Response
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Error(err)
	}
}

func TestErrorResponse(t *testing.T) {
	errResp := suggester.Error{
		Msg: "an error",
	}

	assert.NotEmpty(t, errResp.Error())
}

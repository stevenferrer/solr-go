package helios

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalQueryResponse(t *testing.T) {
	assert := assert.New(t)

	b, err := ioutil.ReadFile("sample-response.json")
	assert.NoError(err)

	// FIXME: assert the fields
	var queryResponse QueryResponse
	err = json.Unmarshal(b, &queryResponse)
	assert.NoError(err)
}

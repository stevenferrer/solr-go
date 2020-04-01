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
	var response SelectResponse
	err = json.Unmarshal(b, &response)
	assert.NoError(err)
}

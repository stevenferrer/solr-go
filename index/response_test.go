package index_test

import (
	"testing"

	"github.com/sf9v/solr-go/index"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	errResp := index.Error{
		Msg: "an error",
	}

	assert.NotEmpty(t, errResp.Error())
}

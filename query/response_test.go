package query_test

import (
	"testing"

	"github.com/sf9v/solr-go/query"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	errResp := query.Error{
		Msg: "an error",
	}

	assert.NotEmpty(t, errResp.Error())
}

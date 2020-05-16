package query_test

import (
	"testing"

	"github.com/stevenferrer/helios/query"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	errResp := query.Error{
		Msg: "an error",
	}

	assert.NotEmpty(t, errResp.Error())
}

package index_test

import (
	"testing"

	"github.com/stevenferrer/helios/index"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	errResp := index.Error{
		Msg: "an error",
	}

	assert.NotEmpty(t, errResp.Error())
}

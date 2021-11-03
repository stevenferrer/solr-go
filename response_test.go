package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stevenferrer/solr-go"
)

func TestResponseError(t *testing.T) {
	err := solr.ResponseError{Msg: "an error"}
	assert.Equal(t, "an error", err.Error())
}

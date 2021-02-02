package solr_test

import (
	"testing"

	"github.com/sf9v/solr-go"
	"github.com/stretchr/testify/assert"
)

func TestResponseError(t *testing.T) {
	err := solr.Error{
		Msg: "an error",
	}

	assert.Equal(t, "an error", err.Error())
}

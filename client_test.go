package helios_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/stevenferrer/helios"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	const (
		host = "192.168.100.19"
		port = 8983
	)

	client := helios.NewClient(host, port)
	b, err := json.Marshal(helios.SimpleQueryRequest{
		Query: "*:*",
	})
	assert.NoError(err)

	response, err := client.Query(ctx, "techproduts", b)
	assert.NoError(err)

	spew.Dump(response)
}

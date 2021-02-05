package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestBuildCollectionParams(t *testing.T) {
	got := solr.NewCollectionParams().
		Name("mycollection").
		NumShards(1).
		ReplicationFactor(1).
		Async("1234").
		BuildParam()

	expect := "async=1234&name=mycollection&numShards=1&replicationFactor=1"
	assert.Equal(t, expect, got)
}

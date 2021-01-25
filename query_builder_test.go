package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestQueryBuilder(t *testing.T) {
	a := assert.New(t)
	qp := solr.NewDisMaxQueryParser("solr rocks")
	qb := solr.NewQueryBuilder(qp)
	got, err := qb.Build()
	a.NoError(err)
	expect := solr.M{"query": "{!dismax}solr rocks"}
	a.Equal(expect, got)
}

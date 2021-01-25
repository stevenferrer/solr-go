package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/solr-go"
)

func TestQueryParsers(t *testing.T) {
	t.Run("standard query parser", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		got, err := solr.NewStandardQueryParser("").BuildQuery()
		a.Error(err)
		a.Empty(got)

		got, err = solr.NewStandardQueryParser("solr rocks").
			WithDf("text").WithOp("AND").Sow().BuildQuery()
		r.NoError(err)
		expect := "{!lucene df='text' q.op='AND' sow=true}solr rocks"
		a.Equal(expect, got)

		got, err = solr.NewStandardQueryParser("").
			WithQ("solr rocks").BuildQuery()
		expect = "{!lucene}solr rocks"
		a.Equal(expect, got)
	})

	t.Run("dismax query parser", func(t *testing.T) {
		a := assert.New(t)

		got, err := solr.NewDisMaxQueryParser("").BuildQuery()
		a.Error(err)
		a.Empty(got)

		got, err = solr.NewDisMaxQueryParser("solr rocks").
			WithAlt("*:*").
			WithQf("one^2.3 two three^0.4").
			WithMm("75%").
			WithPf("one^2.3 two three^0.4").
			WithPs("1").
			WithQs("1").
			WithTie("0.1").
			WithBq("category:food^10").
			WithBf("div(1,sum(1,price))^1.5").
			BuildQuery()
		a.NoError(err)
		expect := "{!dismax q.alt='*:*' qf='one^2.3 two three^0.4' mm='75%' qf='one^2.3 two three^0.4' ps='1' qs='1' tie='0.1' bq='category:food^10' bf='div(1,sum(1,price))^1.5'}solr rocks"
		a.Equal(expect, got)

		got, err = solr.NewDisMaxQueryParser("").
			WithQ("solr rocks").BuildQuery()
		a.NoError(err)
		expect = "{!dismax}solr rocks"
		a.Equal(expect, got)
	})
}

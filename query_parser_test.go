package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/public-safety-canada/solr-go"
)

func TestQueryParsers(t *testing.T) {
	t.Run("standard query parser", func(t *testing.T) {
		a := assert.New(t)

		got := solr.NewStandardQueryParser().
			Tag("certain").BuildParser()
		a.Equal("{!lucene tag=certain}", got)

		got = solr.NewStandardQueryParser().
		    Query("'solr rocks'").
			Df("text").
			Op("AND").
			Sow().
			Rows("100").
			Start("0").
			Fl("id").
			Fq([]string{"(category((11927)))","(category((11838)))"}).
			BuildParser()
		expect := "{!lucene df=text q.op=AND sow=true v='solr rocks' rows=100 start=0 fl=id fq=(category((11927))) fq=(category((11838)))}"
		a.Equal(expect, got)

		got = solr.NewStandardQueryParser().
			Query("'solr rocks'").
			BuildParser()
		expect = "{!lucene v='solr rocks'}"
		a.Equal(expect, got)

		got = solr.NewStandardQueryParser().
		    Query("'solr rocks'").
			Df("text").
			Op("AND").
			Sow().
			Rows("100").
			BuildParser()
		expect = "{!lucene df=text q.op=AND sow=true v='solr rocks' rows=100}"
		a.Equal(expect, got)
	})

	t.Run("dismax query parser", func(t *testing.T) {
		a := assert.New(t)

		got := solr.NewDisMaxQueryParser().BuildParser()
		a.Equal("{!dismax}", got)

		got = solr.NewDisMaxQueryParser().
			Query("'solr rocks'").
			Alt("*:*").
			Qf("'one^2.3 two three^0.4'").
			Mm("75%").
			Pf("'one^2.3 two three^0.4'").
			Ps("1").
			Qs("1").
			Tie("0.1").
			Bq("category:food^10").
			Bf("div(1,sum(1,price))^1.5").
			Rows("100").
			Df("text").
			Op("AND").
			BuildParser()
		expect := `{!dismax q.alt=*:* qf='one^2.3 two three^0.4' mm=75% pf='one^2.3 two three^0.4' ps=1 qs=1 tie=0.1 bq=category:food^10 bf=div(1,sum(1,price))^1.5 v='solr rocks' rows=100 df=text q.op=AND}`
		a.Equal(expect, got)

		got = solr.NewDisMaxQueryParser().
			Query("'solr rocks'").BuildParser()
		expect = "{!dismax v='solr rocks'}"
		a.Equal(expect, got)
	})

	t.Run("extended dismax query parser", func(t *testing.T) {
		a := assert.New(t)

		got := solr.NewExtendedDisMaxQueryParser().BuildParser()
		a.Equal("{!edismax}", got)

		got = solr.NewExtendedDisMaxQueryParser().
			Query("'solr rocks'").
			Alt("*:*").
			Qf("'one^2.3 two three^0.4'").
			Mm("75%").
			Autorelax().
			Pf("'one^2.3 two three^0.4'").
			Ps("1").
			Qs("1").
			Tie("0.1").
			Bq("category:food^10").
			Bf("div(1,sum(1,price))^1.5").
			Uf("title").
			Stopwords("stuff").
			Sow().
			Boost("div(1,sum(1,price))").
			Rows("100").
			Df("text").
			Op("AND").
			Start("0").
			Fl("id").
			Fq([]string{"(category((11927)))","(category((11838)))"}).
			BuildParser()
		expect := `{!edismax q.alt=*:* qf='one^2.3 two three^0.4' mm=75% mm.autorelax=true pf='one^2.3 two three^0.4' ps=1 qs=1 tie=0.1 bq=category:food^10 bf=div(1,sum(1,price))^1.5 uf=title stopwords=stuff sow=true boost=div(1,sum(1,price)) q='solr rocks' rows=100 df=text q.op=AND start=0 fl=id fq=(category((11927))) fq=(category((11838)))}`
		a.Equal(expect, got)

		got = solr.NewExtendedDisMaxQueryParser().
			Query("'solr rocks'").BuildParser()
		expect = "{!edismax q='solr rocks'}"
		a.Equal(expect, got)
	})

	t.Run("parent query parser", func(t *testing.T) {
		a := assert.New(t)
		got := solr.NewParentQueryParser().
			Query("comment:SolrCloud").
			Which("content_type:parent").
			Filters("$childfq").
			ExcludeTags("certain").
			Score("total").
			Tag("top").
			BuildParser()
		expect := `{!parent which=content_type:parent tag=top filters=$childfq excludeTags=certain score=total v=comment:SolrCloud}`
		a.Equal(expect, got)
	})

	t.Run("parent query parser", func(t *testing.T) {
		a := assert.New(t)
		got := solr.NewParentQueryParser().
			Query("comment:SolrCloud").
			Which("content_type:parent").
			Filters("$childfq").
			ExcludeTags("certain").
			Score("total").
			Tag("top").
			BuildParser()
		expect := `{!parent which=content_type:parent tag=top filters=$childfq excludeTags=certain score=total v=comment:SolrCloud}`
		a.Equal(expect, got)
	})

	t.Run("children query parser", func(t *testing.T) {
		a := assert.New(t)
		got := solr.NewChildrenQueryParser().
			Query("$parent").
			Of("$parent").
			Filters("$someFilters").
			ExcludeTags("certain").
			BuildParser()
		expect := `{!child of=$parent filters=$someFilters excludeTags=certain v=$parent}`
		a.Equal(expect, got)
	})

	t.Run("filters query parser", func(t *testing.T) {
		a := assert.New(t)
		got := solr.NewFiltersQueryParser().
			Query("field:text").
			Param("$fqs").
			ExcludeTags("sample").
			BuildParser()
		expect := `{!filters param=$fqs excludeTags=sample v=field:text}`
		a.Equal(expect, got)
	})

}

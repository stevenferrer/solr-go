package solr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sf9v/solr-go"
)

func TestBuildCoreParams(t *testing.T) {
	got := solr.NewCoreParams("mycore").
		DeleteIndex(true).
		DeleteDataDir(true).
		DeleteInstanceDir(true).
		BuildParams()

	expect := "core=mycore&deleteDataDir=true&deleteIndex=true&deleteInstanceDir=true"
	assert.Equal(t, expect, got)
}

func TestBuildCreateCoreParams(t *testing.T) {
	got := solr.NewCreateCoreParams("mycore").
		InstanceDir("mycore").
		Config("solrconfig.xml").
		DataDir("my-data-dir").
		ConfigSet("_default").
		Schema("managed-schema").
		BuildParams()

	expect := "config=solrconfig.xml&configSet=_default&dataDir=my-data-dir&instanceDir=mycore&name=mycore&schema=managed-schema"
	assert.Equal(t, expect, got)
}

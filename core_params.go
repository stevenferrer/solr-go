package solr

import (
	"net/url"
)

type CoreParams struct {
	core string

	// Unload params
	deleteIndex,
	deleteDataDir,
	deleteInstanceDir bool
}

func NewCoreParams(core string) *CoreParams {
	return &CoreParams{core: core}
}

func (c *CoreParams) BuildParams() string {
	vals := &url.Values{}

	if c.core != "" {
		vals.Add("core", c.core)
	}

	if c.deleteIndex {
		vals.Add("deleteIndex", "true")
	}

	if c.deleteDataDir {
		vals.Add("deleteDataDir", "true")
	}

	if c.deleteInstanceDir {
		vals.Add("deleteInstanceDir", "true")
	}

	return vals.Encode()
}

func (c *CoreParams) DeleteIndex(deleteIndex bool) *CoreParams {
	c.deleteIndex = deleteIndex
	return c
}

func (c *CoreParams) DeleteDataDir(deleteDataDir bool) *CoreParams {
	c.deleteDataDir = deleteDataDir
	return c
}

func (c *CoreParams) DeleteInstanceDir(deleteInstanceDir bool) *CoreParams {
	c.deleteInstanceDir = deleteInstanceDir
	return c
}

type CreateCoreParams struct {
	name        string
	instanceDir string
	config      string
	dataDir     string
	configSet   string
	schema      string
}

func NewCreateCoreParams(name string) *CreateCoreParams {
	return &CreateCoreParams{name: name}
}

func (c *CreateCoreParams) InstanceDir(instanceDir string) *CreateCoreParams {
	c.instanceDir = instanceDir
	return c
}

func (c *CreateCoreParams) Config(config string) *CreateCoreParams {
	c.config = config
	return c
}

func (c *CreateCoreParams) DataDir(dataDir string) *CreateCoreParams {
	c.dataDir = dataDir
	return c
}

func (c *CreateCoreParams) ConfigSet(configSet string) *CreateCoreParams {
	c.configSet = configSet
	return c
}

func (c *CreateCoreParams) Schema(schema string) *CreateCoreParams {
	c.schema = schema
	return c
}

// BuildParams builds the parameters
func (c *CreateCoreParams) BuildParams() string {
	vals := &url.Values{}

	if c.name != "" {
		vals.Add("name", c.name)
	}

	if c.instanceDir != "" {
		vals.Add("instanceDir", c.instanceDir)
	}

	if c.config != "" {
		vals.Add("config", c.config)
	}

	if c.dataDir != "" {
		vals.Add("dataDir", c.dataDir)
	}

	if c.configSet != "" {
		vals.Add("configSet", c.configSet)
	}

	if c.schema != "" {
		vals.Add("schema", c.schema)
	}

	return vals.Encode()
}

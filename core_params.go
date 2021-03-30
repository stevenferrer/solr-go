package solr

import (
	"net/url"
)

// CoreParams is the core admin API param builder
type CoreParams struct {
	core string
	deleteIndex,
	deleteDataDir,
	deleteInstanceDir bool
}

// NewCoreParams returns a new CoreParams
func NewCoreParams(core string) *CoreParams {
	return &CoreParams{core: core}
}

// BuildParams builds the parameters
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

// DeleteIndex set to true to remove the index when unloading the core.
// The default is false.
func (c *CoreParams) DeleteIndex(deleteIndex bool) *CoreParams {
	c.deleteIndex = deleteIndex
	return c
}

// DeleteDataDir set to true to remove the data directory and all sub-directories.
// The default is false.
func (c *CoreParams) DeleteDataDir(deleteDataDir bool) *CoreParams {
	c.deleteDataDir = deleteDataDir
	return c
}

// DeleteInstanceDir set to true to remove everything related to the core,
// including the index directory, configuration files and other related files.
// The default is false.
func (c *CoreParams) DeleteInstanceDir(deleteInstanceDir bool) *CoreParams {
	c.deleteInstanceDir = deleteInstanceDir
	return c
}

// CreateCoreParams is the create core param build
type CreateCoreParams struct {
	name        string
	instanceDir string
	config      string
	dataDir     string
	configSet   string
	schema      string
}

// NewCreateCoreParams takes a core name and returns a new CreateCoreParams
func NewCreateCoreParams(name string) *CreateCoreParams {
	return &CreateCoreParams{name: name}
}

// InstanceDir is the directory where files for this core should be stored.
func (c *CreateCoreParams) InstanceDir(instanceDir string) *CreateCoreParams {
	c.instanceDir = instanceDir
	return c
}

// Config is the name of the config file (i.e., solrconfig.xml) relative to instanceDir.
func (c *CreateCoreParams) Config(config string) *CreateCoreParams {
	c.config = config
	return c
}

// Schema is name of the schema file to use for the core.
func (c *CreateCoreParams) Schema(schema string) *CreateCoreParams {
	c.schema = schema
	return c
}

// DataDir is the name of the data directory relative to instanceDir.
func (c *CreateCoreParams) DataDir(dataDir string) *CreateCoreParams {
	c.dataDir = dataDir
	return c
}

// ConfigSet is the name of the configset to use for this core.
func (c *CreateCoreParams) ConfigSet(configSet string) *CreateCoreParams {
	c.configSet = configSet
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

	if c.schema != "" {
		vals.Add("schema", c.schema)
	}

	if c.dataDir != "" {
		vals.Add("dataDir", c.dataDir)
	}

	if c.configSet != "" {
		vals.Add("configSet", c.configSet)
	}

	return vals.Encode()
}

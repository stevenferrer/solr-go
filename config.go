package solr

// CommonProperty is a common property
type CommonProperty struct {
	Name  string
	Value interface{}
}

// ComponentType is a component type
type ComponentType int

// List of component types
const (
	RequestHandler ComponentType = iota
	SearchComponent
	InitParams
	QueryResponseWriter
)

func (ct ComponentType) String() string {
	return [...]string{
		"requesthandler",
		"searchcomponent",
		"initparams",
		"queryresponsewriter",
	}[ct]
}

// Component is a component
type Component struct {
	// ct is the component type
	ct ComponentType
	// name is the component name
	name string
	// class is the component class
	class string
	// m is the component configurations
	m M
}

// NewComponent returns a new Component
func NewComponent(ct ComponentType) *Component {
	return &Component{ct: ct}
}

// Name sets the component name
func (c *Component) Name(name string) *Component {
	c.name = name
	return c
}

// Class sets the component class
func (c *Component) Class(class string) *Component {
	c.class = class
	return c
}

// Config sets the component config
func (c *Component) Config(m M) *Component {
	c.m = m
	return c
}

// BuildComponent builds the component config
func (c *Component) BuildComponent() M {
	m := M{"name": c.name, "class": c.class}

	for k, v := range c.m {
		m[k] = v
	}

	return m
}

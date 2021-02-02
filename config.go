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
	switch ct {
	case RequestHandler:
		return "requesthandler"
	case SearchComponent:
		return "searchcomponent"
	case InitParams:
		return "initparams"
	case QueryResponseWriter:
		return "queryresponsewriter"
	}

	return ""
}

// Component is a component
type Component struct {
	Type   ComponentType
	Name   string
	Class  string
	Values M
}

// UserProperty is a user property
type UserProperty struct{}

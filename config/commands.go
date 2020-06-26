package config

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Commander interface {
	Command() (string, error)
}

// SetPropCommand is a set-property command
type SetPropCommand struct {
	prop string
	val  interface{}
}

// NewSetPropCommand is a factory for SetPropCommand
func NewSetPropCommand(prop string, val interface{}) Commander {
	return SetPropCommand{prop: prop, val: val}
}

// Command implements Commander for SetPropCommand
func (c SetPropCommand) Command() (string, error) {
	m := map[string]interface{}{c.prop: c.val}
	b, err := json.Marshal(m)
	if err != nil {
		return "", errors.Wrap(err, "marshal command")
	}

	return `"set-property": ` + string(b), nil
}

type UnsetPropCommand struct {
	prop string
}

func NewUnsetPropCommand(prop string) Commander {
	return UnsetPropCommand{prop: prop}
}

func (c UnsetPropCommand) Command() (string, error) {
	return fmt.Sprintf(`"unset-property": %q`, c.prop), nil
}

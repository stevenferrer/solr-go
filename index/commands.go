package index

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/stevenferrer/helios"
)

// Cmd is an update command
type Cmd interface {
	// ToCmd returns an update command string
	ToCmd() (string, error)
}

// AddCmd is an add command
type AddCmd struct {
	// CommitWithin option to add the document within the
	// specified number of milliseconds
	CommitWithin int
	// Overwrite indicates if the unique key constraints should be
	// checked to overwrite previous versions of the same document
	Overwrite bool
	// Doc is any arbitrary document to index e.g. map[string]interface{} etc.
	Doc interface{}
}

// ToCmd formats add command
func (c AddCmd) ToCmd() (string, error) {
	cmd := helios.M{}

	if c.CommitWithin > 0 {
		cmd["commitWithin"] = c.CommitWithin
	}

	if c.Overwrite {
		cmd["overwrite"] = c.Overwrite
	}

	cmd["doc"] = c.Doc

	b, err := json.Marshal(cmd)
	if err != nil {
		return "", errors.Wrap(err, "marshal cmd")
	}

	return "\"add\"" + ":" + string(b), nil
}

// DelByQryCmd is a delete by query command
type DelByQryCmd struct {
	Query string
}

// ToCmd formats delete by query command
func (c DelByQryCmd) ToCmd() (string, error) {
	cmd := helios.M{
		"query": c.Query,
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		return "", errors.Wrap(err, "marshal cmd")
	}

	return "\"delete\"" + ":" + string(b), nil
}

// DelByIDsCmd is a delete by list of ids command
type DelByIDsCmd struct {
	IDs []string
}

// ToCmd formats delete by ids command
func (c DelByIDsCmd) ToCmd() (string, error) {
	ids := []string{}
	for _, id := range c.IDs {
		ids = append(ids, fmt.Sprintf("%q", id))
	}

	return "\"delete\"" + ":[" + strings.Join(ids, ",") + "]", nil
}

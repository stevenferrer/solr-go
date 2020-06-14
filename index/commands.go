package index

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Commander is a contract for an update command
type Commander interface {
	// Command builds the command
	Command() (string, error)
}

// AddCommand is an add command
type AddCommand struct {
	// CommitWithin option to add the document within the specified
	// number of milliseconds
	CommitWithin int
	// Overwrite indicates if the unique key constraints should be
	// checked to overwrite previous versions of the same document
	Overwrite bool
	// Doc is any arbitrary document to index e.g. map[string]interface{} etc.
	Doc interface{}
}

// Command formats add command
func (c AddCommand) Command() (string, error) {
	cmd := map[string]interface{}{}

	if c.CommitWithin > 0 {
		cmd["commitWithin"] = c.CommitWithin
	}

	if c.Overwrite {
		cmd["overwrite"] = c.Overwrite
	}

	cmd["doc"] = c.Doc

	b, err := json.Marshal(cmd)
	if err != nil {
		return "", errors.Wrap(err, "marshal command")
	}

	return "\"add\"" + ":" + string(b), nil
}

// DeleteByQueryCommand is a delete-by-query command
type DeleteByQueryCommand struct {
	Query string
}

// Command formats delete by query command
func (c DeleteByQueryCommand) Command() (string, error) {
	cmd := map[string]interface{}{
		"query": c.Query,
	}

	b, err := json.Marshal(cmd)
	if err != nil {
		return "", errors.Wrap(err, "marshal command")
	}

	return "\"delete\"" + ":" + string(b), nil
}

// DeleteByIDsCommand is a delete by list of ids command
type DeleteByIDsCommand struct {
	IDs []string
}

// Command builds the delete-by-ids command
func (c DeleteByIDsCommand) Command() (string, error) {
	ids := []string{}
	for _, id := range c.IDs {
		ids = append(ids, fmt.Sprintf("%q", id))
	}

	return "\"delete\"" + ":[" + strings.Join(ids, ",") + "]", nil
}

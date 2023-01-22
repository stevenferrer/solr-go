package solr

import "fmt"

// M is an alias for map of interface
type M map[string]interface{}

// wrapErr wraps the error with message
func wrapErr(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}

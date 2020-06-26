package config

import (
	"context"
	"errors"
)

// Client is a config API client
type Client interface {
	// GetConfig(ctx context.Context, collection string)
	SendCommands(ctx context.Context, collection string, commands ...Commander) error
}

type client struct{}

// NewClient is a factory for config client
func NewClient() Client { return &client{} }

func (c *client) SendCommands(ctx context.Context, collection string, commands ...Commander) error {
	return errors.New("not implemented")
}

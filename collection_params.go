package solr

import (
	"net/url"
	"strconv"
)

// CollectionParams is the collection API param builder
type CollectionParams struct {
	name              string
	numShards         int
	replicationFactor int
	requestID         string
}

// NewCollectionParams returns a new CollectionParams
func NewCollectionParams() *CollectionParams {
	return &CollectionParams{}
}

// Name sets the collection name
func (c *CollectionParams) Name(name string) *CollectionParams {
	c.name = name
	return c
}

// NumShards sets the number of shards
func (c *CollectionParams) NumShards(ns int) *CollectionParams {
	c.numShards = ns
	return c
}

// ReplicationFactor sets the replication factor
func (c *CollectionParams) ReplicationFactor(rf int) *CollectionParams {
	c.replicationFactor = rf
	return c
}

// Async enable async request with a request ID to track this action
func (c *CollectionParams) Async(requestID string) *CollectionParams {
	c.requestID = requestID
	return c
}

// BuildParam builds the parameters
func (c *CollectionParams) BuildParam() string {
	vals := &url.Values{}

	if c.name != "" {
		vals.Add("name", c.name)
	}

	if c.numShards > 0 {
		vals.Add("numShards", strconv.Itoa(c.numShards))
	}

	if c.replicationFactor > 0 {
		vals.Add("replicationFactor", strconv.Itoa(c.replicationFactor))
	}

	if c.requestID != "" {
		vals.Add("async", c.requestID)
	}

	return vals.Encode()
}

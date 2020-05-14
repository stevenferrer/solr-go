package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/stevenferrer/helios"
)

// Client is the contract for interacting with Solr schema API
type Client interface {
	//  AddField adds a new field definition to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field
	AddField(ctx context.Context, collection string, m helios.M) error
	// DeleteField removes a field definition from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field
	DeleteField(ctx context.Context, collection string, m helios.M) error
	// ReplaceField replaces a field’s definition
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field
	ReplaceField(ctx context.Context, collection string, m helios.M) error
	// AddDynamicField adds a new dynamic field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	AddDynamicField(ctx context.Context, collection string, m helios.M) error
	// DeleteDynamicField deletes a dynamic field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	DeleteDynamicField(ctx context.Context, collection string, m helios.M) error
	// ReplaceDynamicField replaces a dynamic field rule in your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-dynamic-field-rule
	ReplaceDynamicField(ctx context.Context, collection string, m helios.M) error
	// AddCopyField adds a new copy field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-copy-field-rule
	AddCopyField(ctx context.Context, collection string, m helios.M) error
	// DeleteCopyField deletes a copy field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-copy-field-rule
	DeleteCopyField(ctx context.Context, collection string, m helios.M) error
	AddFieldType(ctx context.Context, collection string, m helios.M) error
	DeleteFieldType(ctx context.Context, collection string, m helios.M) error
	ReplaceFieldType(ctx context.Context, collection string, m helios.M) error
}

type client struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewClient is a factory for schema API client
func NewClient(host string, port int, httpClient *http.Client) Client {
	useHTTPS := false
	proto := "http"
	if useHTTPS {
		proto = "https"
	}

	return client{host: host, port: port,
		proto: proto, httpClient: httpClient,
	}
}

func (c client) AddField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "add-field", m)
}

func (c client) DeleteField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "delete-field", m)
}

func (c client) ReplaceField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "replace-field", m)
}

func (c client) AddDynamicField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "add-dynamic-field", m)
}

func (c client) DeleteDynamicField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "delete-dynamic-field", m)
}

func (c client) ReplaceDynamicField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "replace-dynamic-field", m)
}

func (c client) AddCopyField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "add-copy-field", m)
}

func (c client) DeleteCopyField(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "delete-copy-field", m)
}

func (c client) AddFieldType(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "add-field-type", m)
}

func (c client) DeleteFieldType(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "delete-field-type", m)
}
func (c client) ReplaceFieldType(ctx context.Context, collection string, m helios.M) error {
	return c.doCmd(ctx, collection, "replace-field-type", m)
}

func (c client) doCmd(ctx context.Context, collection, op string, m helios.M) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(helios.M{op: m})
	if err != nil {
		return errors.Wrap(err, "marshal request")
	}

	var httpReq *http.Request
	httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, theURL.String(), bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "http do request")
	}

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return resp.Error
	}

	return nil
}

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
	AddField(ctx context.Context, coll string, fld Field) error
	// DeleteField removes a field definition from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field
	DeleteField(ctx context.Context, coll string, fld Field) error
	// ReplaceField replaces a field’s definition
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field
	ReplaceField(ctx context.Context, coll string, fld Field) error
	// AddDynamicField adds a new dynamic field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	AddDynamicField(ctx context.Context, coll string, fld Field) error
	// DeleteDynamicField deletes a dynamic field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	DeleteDynamicField(ctx context.Context, coll string, fld Field) error
	// ReplaceDynamicField replaces a dynamic field rule in your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-dynamic-field-rule
	ReplaceDynamicField(ctx context.Context, coll string, fld Field) error
	// AddFieldType adds a new field type to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field-type
	AddFieldType(ctx context.Context, coll string, fldTyp FieldType) error
	// DeleteFieldType removes a field type from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field-type
	DeleteFieldType(ctx context.Context, coll string, fldTyp FieldType) error
	// ReplaceFieldType replaces a field type in your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field-type
	ReplaceFieldType(ctx context.Context, coll string, fldTyp FieldType) error
	// AddCopyField adds a new copy field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-copy-field-rule
	AddCopyField(ctx context.Context, coll string, cpyFld CopyField) error
	// DeleteCopyField deletes a copy field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-copy-field-rule
	DeleteCopyField(ctx context.Context, coll string, cpyFld CopyField) error
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

func (c client) AddField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "add-field", fld)
}

func (c client) DeleteField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "delete-field", fld)
}

func (c client) ReplaceField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "replace-field", fld)
}

func (c client) AddDynamicField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "add-dynamic-field", fld)
}

func (c client) DeleteDynamicField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "delete-dynamic-field", fld)
}

func (c client) ReplaceDynamicField(ctx context.Context, coll string, fld Field) error {
	return c.doCmd(ctx, coll, "replace-dynamic-field", fld)
}

func (c client) AddCopyField(ctx context.Context, coll string, cpyFld CopyField) error {
	return c.doCmd(ctx, coll, "add-copy-field", cpyFld)
}

func (c client) DeleteCopyField(ctx context.Context, coll string, cpyFld CopyField) error {
	return c.doCmd(ctx, coll, "delete-copy-field", cpyFld)
}

func (c client) AddFieldType(ctx context.Context, coll string, fldTyp FieldType) error {
	return c.doCmd(ctx, coll, "add-field-type", fldTyp)
}

func (c client) DeleteFieldType(ctx context.Context, coll string, fldTyp FieldType) error {
	return c.doCmd(ctx, coll, "delete-field-type", fldTyp)
}
func (c client) ReplaceFieldType(ctx context.Context, coll string, fldTyp FieldType) error {
	return c.doCmd(ctx, coll, "replace-field-type", fldTyp)
}

func (c client) doCmd(ctx context.Context, coll, cmd string, body interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema",
		c.proto, c.host, c.port, coll))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(helios.M{cmd: body})
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

package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/stevenferrer/solr-go/types"
)

// Client is the contract for interacting with Solr schema API
type Client interface {
	// GetSchema returns the schema details
	GetSchema(ctx context.Context, collection string) (*Schema, error)
	// ListFields returns the list of fields
	ListFields(ctx context.Context, collection string) ([]Field, error)
	// GetField returns a field details
	GetField(ctx context.Context, collection, fldNm string) (*Field, error)
	// ListDynamicFields returns list of dynamic fields
	ListDynamicFields(ctx context.Context, collection string) ([]Field, error)
	// ListFieldTypes returns the list of field types
	ListFieldTypes(ctx context.Context, collection string) ([]FieldType, error)
	// ListCopyFields returns the list of copy fields
	ListCopyFields(ctx context.Context, collection string) ([]CopyField, error)

	//  AddField adds a new field definition to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field
	AddField(ctx context.Context, collection string, fld Field) error
	// DeleteField removes a field definition from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field
	DeleteField(ctx context.Context, collection string, fld Field) error
	// ReplaceField replaces a field’s definition
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field
	ReplaceField(ctx context.Context, collection string, fld Field) error
	// AddDynamicField adds a new dynamic field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	AddDynamicField(ctx context.Context, collection string, fld Field) error
	// DeleteDynamicField deletes a dynamic field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-dynamic-field-rule
	DeleteDynamicField(ctx context.Context, collection string, fld Field) error
	// ReplaceDynamicField replaces a dynamic field rule in your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-dynamic-field-rule
	ReplaceDynamicField(ctx context.Context, collection string, fld Field) error
	// AddFieldType adds a new field type to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-field-type
	AddFieldType(ctx context.Context, collection string, fldTyp FieldType) error
	// DeleteFieldType removes a field type from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-field-type
	DeleteFieldType(ctx context.Context, collection string, fldTyp FieldType) error
	// ReplaceFieldType replaces a field type in your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#replace-a-field-type
	ReplaceFieldType(ctx context.Context, collection string, fldTyp FieldType) error
	// AddCopyField adds a new copy field rule to your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#add-a-new-copy-field-rule
	AddCopyField(ctx context.Context, collection string, cpyFld CopyField) error
	// DeleteCopyField deletes a copy field rule from your schema
	// Reference: https://lucene.apache.org/solr/guide/8_5/schema-api.html#delete-a-copy-field-rule
	DeleteCopyField(ctx context.Context, collection string, cpyFld CopyField) error
}

type client struct {
	host       string
	port       int
	proto      string
	httpClient *http.Client
}

// NewClient is a factory for schema API client
func NewClient(host string, port int) Client {
	proto := "http"
	return &client{
		host:  host,
		port:  port,
		proto: proto,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// NewClientWithHTTPClient is a factory for schema API client with custom http client
func NewClientWithHTTPClient(host string, port int, httpClient *http.Client) Client {
	proto := "http"
	return &client{host: host, port: port,
		proto: proto, httpClient: httpClient,
	}
}

func (c client) GetSchema(ctx context.Context, collection string) (*Schema, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.Schema, nil
}

func (c client) ListFields(ctx context.Context, collection string) ([]Field, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema/fields",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.Fields, nil
}

func (c client) ListDynamicFields(ctx context.Context, collection string) ([]Field, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema/dynamicfields",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.DynamicFields, nil
}

func (c client) GetField(ctx context.Context, collection, fldNm string) (*Field, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema/fields/%s",
		c.proto, c.host, c.port, collection, fldNm))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.Field, nil
}

func (c client) ListFieldTypes(ctx context.Context, collection string) ([]FieldType, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema/fieldtypes",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.FieldTypes, nil
}

func (c client) ListCopyFields(ctx context.Context, collection string) ([]CopyField, error) {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema/copyfields",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}

	var resp *Response
	resp, err = c.doRtrv(ctx, theURL.String())
	if err != nil {
		return nil, err
	}
	return resp.CopyFields, nil
}

func (c client) AddField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "add-field", fld)
}

func (c client) DeleteField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "delete-field", fld)
}

func (c client) ReplaceField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "replace-field", fld)
}

func (c client) AddDynamicField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "add-dynamic-field", fld)
}

func (c client) DeleteDynamicField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "delete-dynamic-field", fld)
}

func (c client) ReplaceDynamicField(ctx context.Context, collection string, fld Field) error {
	return c.doMdfy(ctx, collection, "replace-dynamic-field", fld)
}

func (c client) AddCopyField(ctx context.Context, collection string, cpyFld CopyField) error {
	return c.doMdfy(ctx, collection, "add-copy-field", cpyFld)
}

func (c client) DeleteCopyField(ctx context.Context, collection string, cpyFld CopyField) error {
	return c.doMdfy(ctx, collection, "delete-copy-field", cpyFld)
}

func (c client) AddFieldType(ctx context.Context, collection string, fldTyp FieldType) error {
	return c.doMdfy(ctx, collection, "add-field-type", fldTyp)
}

func (c client) DeleteFieldType(ctx context.Context, collection string, fldTyp FieldType) error {
	return c.doMdfy(ctx, collection, "delete-field-type", fldTyp)
}
func (c client) ReplaceFieldType(ctx context.Context, collection string, fldTyp FieldType) error {
	return c.doMdfy(ctx, collection, "replace-field-type", fldTyp)
}

func (c client) doMdfy(ctx context.Context, collection, cmd string, body interface{}) error {
	theURL, err := url.Parse(fmt.Sprintf("%s://%s:%d/solr/%s/schema",
		c.proto, c.host, c.port, collection))
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	var b []byte
	b, err = json.Marshal(types.M{cmd: body})
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

func (c client) doRtrv(ctx context.Context, theURL string) (*Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, theURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", "application/json")

	var httpResp *http.Response
	httpResp, err = c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "http do request")
	}

	var resp Response
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	if httpResp.StatusCode > http.StatusOK {
		return nil, resp.Error
	}

	return &resp, nil
}

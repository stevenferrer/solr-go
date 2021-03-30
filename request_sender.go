package solr

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// RequestSender is an HTTP request sender
type RequestSender interface {
	SendRequest(ctx context.Context, method, urlStr,
		contentType string, body io.Reader) (*http.Response, error)
}

type basicAuth struct {
	username, password string
}

// DefaultRequestSender is the default HTTP request sender
type DefaultRequestSender struct {
	httpClient *http.Client
	basicAuth  *basicAuth
}

var _ RequestSender = (*DefaultRequestSender)(nil)

// NewDefaultRequestSender returns a new DefaultRequestSender
func NewDefaultRequestSender() *DefaultRequestSender {
	return &DefaultRequestSender{
		httpClient: http.DefaultClient,
	}
}

// WithHTTPClient overrides the default HTTP client
func (rs *DefaultRequestSender) WithHTTPClient(httpClient *http.Client) *DefaultRequestSender {
	rs.httpClient = httpClient
	return rs
}

// WithBasicAuth sets the basic auth credentials
func (rs *DefaultRequestSender) WithBasicAuth(username, password string) *DefaultRequestSender {
	rs.basicAuth = &basicAuth{username: username, password: password}
	return rs
}

// SendRequest builds and sends the HTTP request
func (rs *DefaultRequestSender) SendRequest(ctx context.Context, httpMethod,
	urlStr, contentType string, body io.Reader) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, httpMethod, urlStr, body)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("content-type", contentType)

	// include basic auth if available
	if rs.basicAuth != nil {
		httpReq.SetBasicAuth(rs.basicAuth.username, rs.basicAuth.password)
	}

	var httpResp *http.Response
	httpResp, err = rs.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
	}

	return httpResp, nil
}

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

// DefaultRequestSender is the default HTTP request sender
type DefaultRequestSender struct {
	httpClient *http.Client
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

// SendRequest builds and sends the HTTP request
func (rs *DefaultRequestSender) SendRequest(
	ctx context.Context, httpMethod,
	urlStr, contentType string,
	body io.Reader,
) (*http.Response, error) {
	httpReq, err := http.NewRequestWithContext(ctx, httpMethod, urlStr, body)
	if err != nil {
		return nil, errors.Wrap(err, "new http request")
	}
	httpReq.Header.Add("Content-Type", contentType)

	var httpResp *http.Response
	httpResp, err = rs.httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "send http request")
	}

	return httpResp, nil
}

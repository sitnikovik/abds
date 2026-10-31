// Package httpclient provides a fake HTTP client for tests.
package httpclient

import (
	"io"
	"net/http"
	"strings"
)

// Response is one canned response of the fake, or an error instead of it.
type Response struct {
	// status is the status code of the response.
	status int
	// body is the body of the response.
	body string
	// err is returned instead of a response when it is not nil.
	err error
}

// ResponseOption configures a canned Response.
type ResponseOption func(*Response)

// WithStatus sets the status code of the response.
func WithStatus(status int) ResponseOption {
	return func(r *Response) {
		r.status = status
	}
}

// WithBody sets the body of the response.
func WithBody(body string) ResponseOption {
	return func(r *Response) {
		r.body = body
	}
}

// WithErr makes the fake return err instead of the response.
func WithErr(err error) ResponseOption {
	return func(r *Response) {
		r.err = err
	}
}

// NewResponse creates a canned Response with status 200 and an empty body, the options change it.
func NewResponse(opts ...ResponseOption) Response {
	r := Response{
		status: http.StatusOK,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// Client is a fake HTTP client that answers requests with predefined responses in order.
//
// Once every response is used, the last one is returned for the remaining requests.
type Client struct {
	// responses are returned in order.
	responses []Response
	// calls is the number of requests received so far.
	calls int
}

// NewClient creates a fake Client that answers every request with status and body, or fails with err.
func NewClient(status int, body string, err error) *Client {
	return NewSequence(
		NewResponse(
			WithStatus(status),
			WithBody(body),
			WithErr(err),
		),
	)
}

// NewSequence creates a fake Client that answers requests with responses in order.
func NewSequence(responses ...Response) *Client {
	return &Client{responses: responses}
}

// Get returns the next predefined response.
func (c *Client) Get(_ string) (*http.Response, error) {
	return c.Do(nil)
}

// Do returns the next predefined response, or the error it holds.
//
// The response declares a JSON content type.
func (c *Client) Do(_ *http.Request) (*http.Response, error) {
	i := c.calls
	if i >= len(c.responses) {
		i = len(c.responses) - 1
	}
	c.calls++
	canned := c.responses[i]
	if canned.err != nil {
		return nil, canned.err
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: canned.status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(canned.body)),
	}, nil
}

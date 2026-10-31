// Package httpclient provides a fake HTTP client for tests.
package httpclient

import (
	"io"
	"net/http"
	"strings"
)

// Canned is a response the fake returns for one request, or an error instead of it.
type Canned struct {
	// Status is the status code of the response.
	Status int
	// Body is the body of the response.
	Body string
	// Err is returned instead of a response when it is not nil.
	Err error
}

// Client is a fake HTTP client that answers requests with predefined responses in order.
//
// Once every response is used, the last one is returned for the remaining requests.
type Client struct {
	// responses are returned in order.
	responses []Canned
	// calls is the number of requests received so far.
	calls int
}

// NewClient creates a fake Client that answers every request with status and body, or fails with err.
func NewClient(status int, body string, err error) *Client {
	return NewSequence(Canned{Status: status, Body: body, Err: err})
}

// NewSequence creates a fake Client that answers requests with responses in order.
func NewSequence(responses ...Canned) *Client {
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
	if canned.Err != nil {
		return nil, canned.Err
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: canned.Status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(canned.Body)),
	}, nil
}

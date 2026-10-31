// Package httpx builds the HTTP requests the autograder sends and checks the responses it gets.
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/errs"
)

// Doer sends HTTP requests and returns the responses.
type Doer interface {
	// Do sends req and returns the response.
	Do(req *http.Request) (*http.Response, error)
}

// NewRequest creates a request with method and url.
//
// When body is not nil it is encoded as JSON and sent as the request body.
func NewRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode body: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// ReadBody reads the whole body of rsp, the caller closes the body.
func ReadBody(rsp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return data, nil
}

// Status returns nil if rsp has the status code want.
func Status(rsp *http.Response, want int) error {
	return errs.Wrap(
		number.Equal(rsp.StatusCode, want),
		fmt.Sprintf("status code must be %d", want),
	)
}

// JSONContent returns nil if rsp declares a JSON content type.
func JSONContent(rsp *http.Response) error {
	return errs.Wrap(
		strings.Contains(rsp.Header.Get("Content-Type"), "application/json"),
		"content type must be application/json",
	)
}

// Decode unmarshals the JSON body into dst.
func Decode(body []byte, dst any) error {
	if err := json.Unmarshal(body, dst); err != nil {
		return fmt.Errorf("decode body: %w", err)
	}
	return nil
}

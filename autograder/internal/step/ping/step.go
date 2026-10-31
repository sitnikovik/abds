// Package ping provides the step that checks the health endpoint of the application.
package ping

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// httpClient is the part of the HTTP client the step uses.
type httpClient interface {
	// Get sends a GET request to url.
	Get(url string) (*http.Response, error)
}

// Step checks that the application answers GET /ping with 200 and an empty body.
type Step struct {
	// client sends the request.
	client httpClient
	// url is the full URL of the ping endpoint.
	url string
}

// NewStep creates a Step that sends its request to url through client.
func NewStep(client httpClient, url string) Step {
	return Step{
		client: client,
		url:    url,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "ping"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /ping returns 200 OK with an empty body."
}

// Run sends GET to the ping endpoint and checks the status code and the body.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Step) Run(_ context.Context, _ step.Variables) error {
	rsp, err := s.client.Get(s.url)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	return expect.All(
		number.Equal(rsp.StatusCode, http.StatusOK),
		strings.Empty(string(body)),
	)
}

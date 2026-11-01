// Package ping provides the step that checks the health endpoint of the application.
package ping

import (
	"net/http"
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

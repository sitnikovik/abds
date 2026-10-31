// Package badrequest provides the step that checks the tariff endpoint without a date.
package badrequest

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// Step checks that GET /tariffs without a date answers 400 with a message.
type Step struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that sends its request to baseURL through cli.
func NewStep(cli httpx.Doer, baseURL string) Step {
	return Step{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "tariffs bad request"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /tariffs without a date answers 400 with a message."
}

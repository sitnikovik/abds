// Package notfound provides the step that checks the owner endpoint for an unknown owner.
package notfound

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// missingUserID is an id that no owner has.
const missingUserID = "999999999"

// Step checks that GET /users/{id} answers 404 with a message for an unknown owner.
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
	return "users get not found"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /users/{id} answers 404 with a message for an unknown owner."
}

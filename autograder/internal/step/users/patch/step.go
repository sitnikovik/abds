// Package patch provides the step that checks the owner's phone update.
package patch

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// patchPhone is the phone the step sets on the owner.
const patchPhone = "+79999999999"

// Step checks that PATCH /users/{id} changes the phone of the owner and answers 204 without a body.
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
	return "users patch"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that PATCH /users/{id} updates the phone and answers 204 No Content."
}

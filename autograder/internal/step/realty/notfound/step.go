// Package notfound provides the step that checks the object endpoint for an unknown object.
package notfound

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// missingRealtyID is an id that no object has.
const missingRealtyID = "999999999"

// Step checks that GET /realty/{id} answers 404 with a message for an unknown object.
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
	return "realty get not found"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /realty/{id} answers 404 with a message for an unknown object."
}

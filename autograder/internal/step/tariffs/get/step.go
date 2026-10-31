// Package get provides the step that checks the tariff endpoint.
package get

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// lookupDate is the date the tariff in force is requested for.
const lookupDate = "2026-10-27"

// Step checks that GET /tariffs?date= returns the tariff version in force on the date.
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
	return "tariffs get"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /tariffs?date= returns the tariff version in force on the date."
}

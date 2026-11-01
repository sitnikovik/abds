// Package list provides the step that checks the page of invoices of an object.
package list

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// pageLimit is the page size the step requests.
const pageLimit = 10

// Step checks that GET /realty/{id}/invoices returns a page of invoices with the total count, newest first.
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
	return "invoices list"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that GET /realty/{id}/invoices returns a page of invoices, the total count and newest first order."
}

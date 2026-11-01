// Package pagination provides the step that checks the pagination of the invoice endpoint.
package pagination

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// defaultPageLimit is the page size the API uses when the limit is not passed.
const defaultPageLimit = 10

// Step checks the default page size of GET /realty/{id}/invoices and that offset moves to the next page.
type Step struct {
	// cli sends the requests.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that sends its requests to baseURL through cli.
func NewStep(cli httpx.Doer, baseURL string) Step {
	return Step{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "invoices pagination"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks the default page size of invoices and that offset moves to the next page."
}

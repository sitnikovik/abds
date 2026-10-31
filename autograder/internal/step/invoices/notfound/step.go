// Package notfound provides the step that checks the invoice endpoint for an unknown object.
package notfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that GET /realty/{id}/invoices for an unknown object answers 404 with a message.
type Step struct {
	// Step sends the request and checks the response.
	endpoint.Step
}

// NewStep creates a Step that sends its request to baseURL through cli.
func NewStep(cli httpx.Doer, baseURL string) Step {
	return Step{
		Step: endpoint.NewStep(
			cli,
			baseURL,
			step.NewDesc(
				"invoices not found",
				"Checks that GET /realty/{id}/invoices for an unknown object answers 404 Not Found with a message.",
			),
			endpoint.WithPath("/realty/999999999/invoices"),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

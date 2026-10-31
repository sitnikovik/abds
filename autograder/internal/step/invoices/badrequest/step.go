// Package badrequest provides the step that checks the invoice endpoint with a negative limit.
package badrequest

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Step checks that GET /realty/{id}/invoices with a negative limit answers 400 with a message.
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
				"invoices bad request",
				"Checks that GET /realty/{id}/invoices with a negative limit answers 400 Bad Request with a message.",
			),
			endpoint.WithPath("/realty/%s/invoices?limit=-1&offset=0", fixture.RealtyID),
			endpoint.WithStatus(http.StatusBadRequest),
		),
	}
}

// Package notfound provides the step that checks the tariff endpoint before every version.
package notfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that GET /tariffs for a date before every version answers 404 with a message.
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
				"tariffs not found",
				"Checks that GET /tariffs for a date before every version answers 404 Not Found with a message.",
			),
			endpoint.WithPath("/tariffs?date=1900-01-01"),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

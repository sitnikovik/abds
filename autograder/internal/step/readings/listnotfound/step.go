// Package listnotfound provides the step that checks the readings endpoint for an unknown object.
package listnotfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that GET /realty/{id}/readings for an unknown object answers 404 with a message.
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
				"readings list not found",
				"Checks that GET /realty/{id}/readings for an unknown object answers 404 Not Found with a message.",
			),
			endpoint.WithPath("/realty/999999999/readings?date_from=2026-10-01&date_to=2026-10-31"),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

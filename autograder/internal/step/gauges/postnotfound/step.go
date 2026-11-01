// Package postnotfound provides the step that checks the gauge installation for an unknown object.
package postnotfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that POST /realty/{id}/gauges for an unknown object answers 404 with a message.
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
				"gauges post not found",
				"Checks that POST /realty/{id}/gauges for an unknown object answers 404 Not Found with a message.",
			),
			endpoint.WithMethod(http.MethodPost),
			endpoint.WithPath("/realty/999999999/gauges"),
			endpoint.WithBody(map[string]int{"gauge_id": 2}),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

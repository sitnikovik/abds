// Package invaliddate provides the step that checks the tariff endpoint with a malformed date.
package invaliddate

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that GET /tariffs with a malformed date answers 400 with a message.
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
				"tariffs invalid date",
				"Checks that GET /tariffs with a malformed date answers 400 Bad Request with a message.",
			),
			endpoint.WithPath("/tariffs?date=not-a-date"),
			endpoint.WithStatus(http.StatusBadRequest),
		),
	}
}

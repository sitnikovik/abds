// Package postbadrequest provides the step that checks the reading submission with a malformed date.
package postbadrequest

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Step checks that POST /realty/{id}/readings with a malformed date answers 400 with a message.
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
				"readings post bad request",
				"Checks that POST /realty/{id}/readings with a malformed date answers 400 Bad Request with a message.",
			),
			endpoint.WithMethod(http.MethodPost),
			endpoint.WithPath("/realty/%s/readings", fixture.RealtyID),
			endpoint.WithBody(map[string]string{"reading_date": "not-a-date"}),
			endpoint.WithStatus(http.StatusBadRequest),
		),
	}
}

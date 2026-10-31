// Package patchbadrequest provides the step that checks the owner update without fields.
package patchbadrequest

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Step checks that PATCH /users/{id} without fields answers 400 with a message.
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
				"users patch bad request",
				"Checks that PATCH /users/{id} without fields answers 400 Bad Request with a message.",
			),
			endpoint.WithMethod(http.MethodPatch),
			endpoint.WithPath("/users/%s", fixture.UserID),
			endpoint.WithBody(map[string]string{}),
			endpoint.WithStatus(http.StatusBadRequest),
		),
	}
}

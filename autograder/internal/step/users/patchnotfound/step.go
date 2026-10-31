// Package patchnotfound provides the step that checks the owner update for an unknown owner.
package patchnotfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// Step checks that PATCH /users/{id} for an unknown owner answers 404 with a message.
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
				"users patch not found",
				"Checks that PATCH /users/{id} for an unknown owner answers 404 Not Found with a message.",
			),
			endpoint.WithMethod(http.MethodPatch),
			endpoint.WithPath("/users/999999999"),
			endpoint.WithBody(map[string]string{"phone": "+79999999998"}),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

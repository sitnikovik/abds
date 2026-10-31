package users

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// NewPatchBadRequest creates a step that checks PATCH /users/{id} without fields answers 400 with a message.
func NewPatchBadRequest(cli httpx.Doer, baseURL string) endpoint.Step {
	return endpoint.NewStep(
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
	)
}

// NewPatchNotFound creates a step that checks PATCH /users/{id} for an unknown owner answers 404 with a message.
func NewPatchNotFound(cli httpx.Doer, baseURL string) endpoint.Step {
	return endpoint.NewStep(
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
	)
}

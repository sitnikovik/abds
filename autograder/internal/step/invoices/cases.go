package invoices

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/httpcase"
)

// NewBadRequest creates a step that checks GET /realty/{id}/invoices with a negative limit answers 400 with a message.
func NewBadRequest(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"invoices bad request",
			"Checks that GET /realty/{id}/invoices with a negative limit answers 400 Bad Request with a message.",
		),
		httpcase.WithPath("/realty/%s/invoices?limit=-1&offset=0", fixture.RealtyID),
		httpcase.WithStatus(http.StatusBadRequest),
	)
}

// NewNotFound creates a step that checks GET /realty/{id}/invoices for an unknown object answers 404 with a message.
func NewNotFound(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"invoices not found",
			"Checks that GET /realty/{id}/invoices for an unknown object answers 404 Not Found with a message.",
		),
		httpcase.WithPath("/realty/999999999/invoices"),
		httpcase.WithStatus(http.StatusNotFound),
	)
}

package tariffs

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/httpcase"
)

// NewInvalidDate creates a step that checks GET /tariffs with a malformed date answers 400 with a message.
func NewInvalidDate(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"tariffs invalid date",
			"Checks that GET /tariffs with a malformed date answers 400 Bad Request with a message.",
		),
		httpcase.WithPath("/tariffs?date=not-a-date"),
		httpcase.WithStatus(http.StatusBadRequest),
	)
}

// NewNotFound creates a step that checks GET /tariffs for a date before every version answers 404 with a message.
func NewNotFound(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"tariffs not found",
			"Checks that GET /tariffs for a date before every version answers 404 Not Found with a message.",
		),
		httpcase.WithPath("/tariffs?date=1900-01-01"),
		httpcase.WithStatus(http.StatusNotFound),
	)
}

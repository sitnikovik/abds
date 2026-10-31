package readings

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/httpcase"
)

// NewListNotFound creates a step that checks GET /realty/{id}/readings for an unknown object answers 404 with a message.
func NewListNotFound(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"readings list not found",
			"Checks that GET /realty/{id}/readings for an unknown object answers 404 Not Found with a message.",
		),
		httpcase.WithPath("/realty/999999999/readings?date_from=2026-10-01&date_to=2026-10-31"),
		httpcase.WithStatus(http.StatusNotFound),
	)
}

// NewPostBadRequest creates a step that checks POST /realty/{id}/readings with a malformed date answers 400 with a message.
func NewPostBadRequest(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"readings post bad request",
			"Checks that POST /realty/{id}/readings with a malformed date answers 400 Bad Request with a message.",
		),
		httpcase.WithMethod(http.MethodPost),
		httpcase.WithPath("/realty/%s/readings", fixture.RealtyID),
		httpcase.WithBody(map[string]string{"reading_date": "not-a-date"}),
		httpcase.WithStatus(http.StatusBadRequest),
	)
}

// NewPostNotFound creates a step that checks POST /realty/{id}/readings for an unknown object answers 404 with a message.
func NewPostNotFound(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"readings post not found",
			"Checks that POST /realty/{id}/readings for an unknown object answers 404 Not Found with a message.",
		),
		httpcase.WithMethod(http.MethodPost),
		httpcase.WithPath("/realty/999999999/readings"),
		httpcase.WithBody(map[string]any{
			"user_id":      1,
			"gauge_id":     1,
			"tariff_id":    1,
			"reading_date": readingDate,
			"day_kw":       dayKW,
			"night_kw":     nightKW,
		}),
		httpcase.WithStatus(http.StatusNotFound),
	)
}

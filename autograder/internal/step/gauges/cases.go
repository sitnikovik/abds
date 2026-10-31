package gauges

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/httpcase"
)

// NewPostBadRequest creates a step that checks POST /realty/{id}/gauges without a gauge answers 400 with a message.
func NewPostBadRequest(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"gauges post bad request",
			"Checks that POST /realty/{id}/gauges without a gauge answers 400 Bad Request with a message.",
		),
		httpcase.WithMethod(http.MethodPost),
		httpcase.WithPath("/realty/%s/gauges", fixture.RealtyID),
		httpcase.WithBody(map[string]string{}),
		httpcase.WithStatus(http.StatusBadRequest),
	)
}

// NewPostNotFound creates a step that checks POST /realty/{id}/gauges for an unknown object answers 404 with a message.
func NewPostNotFound(cli httpx.Doer, baseURL string) httpcase.Step {
	return httpcase.NewStep(
		cli,
		baseURL,
		step.NewDesc(
			"gauges post not found",
			"Checks that POST /realty/{id}/gauges for an unknown object answers 404 Not Found with a message.",
		),
		httpcase.WithMethod(http.MethodPost),
		httpcase.WithPath("/realty/999999999/gauges"),
		httpcase.WithBody(map[string]int{"gauge_id": 2}),
		httpcase.WithStatus(http.StatusNotFound),
	)
}

// Package postnotfound provides the step that checks the reading submission for an unknown object.
package postnotfound

import (
	"net/http"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
)

// readingDate is the date of the reading the step sends.
const readingDate = "2026-10-01"

// dayKW is the day counter value the step sends.
const dayKW = 1250

// nightKW is the night counter value the step sends.
const nightKW = 300

// Step checks that POST /realty/{id}/readings for an unknown object answers 404 with a message.
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
				"readings post not found",
				"Checks that POST /realty/{id}/readings for an unknown object answers 404 Not Found with a message.",
			),
			endpoint.WithMethod(http.MethodPost),
			endpoint.WithPath("/realty/999999999/readings"),
			endpoint.WithBody(map[string]any{
				"user_id":      1,
				"gauge_id":     1,
				"tariff_id":    1,
				"reading_date": readingDate,
				"day_kw":       dayKW,
				"night_kw":     nightKW,
			}),
			endpoint.WithStatus(http.StatusNotFound),
		),
	}
}

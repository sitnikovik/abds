// Package post provides the step that checks the reading submission.
package post

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// readingDate is the date of the reading the step sends.
const readingDate = "2026-10-01"

// dayKW is the day counter value the step sends.
const dayKW = 1250

// nightKW is the night counter value the step sends.
const nightKW = 300

// postBody is the JSON body of POST /realty/{id}/readings.
type postBody struct {
	// ReadingDate is the date of the reading in the YYYY-MM-DD format.
	ReadingDate string `json:"reading_date"`
	// UserID is the id of the owner who sends the reading.
	UserID uint64 `json:"user_id"`
	// GaugeID is the id of the gauge installed on the object.
	GaugeID uint64 `json:"gauge_id"`
	// TariffID is the id of the tariff the reading is billed by.
	TariffID uint64 `json:"tariff_id"`
	// DayKW is the day counter value.
	DayKW uint32 `json:"day_kw"`
	// NightKW is the night counter value.
	NightKW uint32 `json:"night_kw"`
}

// Step checks that POST /realty/{id}/readings accepts a reading and answers 201 without a body.
type Step struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that sends its request to baseURL through cli.
func NewStep(cli httpx.Doer, baseURL string) Step {
	return Step{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "readings post"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that POST /realty/{id}/readings accepts a reading and answers 201 Created."
}

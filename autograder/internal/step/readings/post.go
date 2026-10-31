// Package readings provides the steps that check the meter readings endpoints.
package readings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
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

// Post checks that POST /realty/{id}/readings accepts a reading and answers 201 without a body.
type Post struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewPost creates a Post step that sends its request to baseURL through cli.
func NewPost(cli httpx.Doer, baseURL string) Post {
	return Post{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Post) Name() string {
	return "readings post"
}

// Description returns a brief explanation of what the step checks.
func (s Post) Description() string {
	return "Checks that POST /realty/{id}/readings accepts a reading and answers 201 Created."
}

// Run sends a reading for the object stored under the realty id and checks the response.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Post) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.Require(fixture.RealtyID)
	if err != nil {
		return err
	}
	userID, err := vars.RequireUint64(fixture.UserID)
	if err != nil {
		return err
	}
	gaugeID, err := vars.RequireUint64(fixture.GaugeID)
	if err != nil {
		return err
	}
	tariffID, err := vars.RequireUint64(fixture.TariffID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPost,
		s.baseURL+"/realty/"+id+"/readings",
		postBody{
			UserID:      userID,
			GaugeID:     gaugeID,
			TariffID:    tariffID,
			ReadingDate: readingDate,
			DayKW:       dayKW,
			NightKW:     nightKW,
		},
	)
	if err != nil {
		return err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return err
	}
	return expect.All(
		httpx.Status(rsp, http.StatusCreated),
		errs.Wrap(
			strings.Empty(string(body)),
			"body must be empty",
		),
	)
}

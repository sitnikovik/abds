package readings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Correction checks that the readings of the date show only the last value sent for it (ReplacingMergeTree).
//
// It must run after Post, which sends the reading of the date.
type Correction struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewCorrection creates a Correction step that sends its request to baseURL through cli.
func NewCorrection(cli httpx.Doer, baseURL string) Correction {
	return Correction{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Correction) Name() string {
	return "readings correction"
}

// Description returns a brief explanation of what the step checks.
func (s Correction) Description() string {
	return "Checks that the readings of a date show one actual value after a correction."
}

// Run requests the readings of the date sent by Post and checks that exactly one row with the sent values is returned.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Correction) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.Require(fixture.RealtyID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/realty/%s/readings?date_from=%s&date_to=%s", s.baseURL, id, readingDate, readingDate),
		nil,
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
	if err := expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp)); err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	var got struct {
		Readings []struct {
			DayKW   uint32 `json:"day_kw"`
			NightKW uint32 `json:"night_kw"`
		} `json:"readings"`
	}
	if err := httpx.Decode(body, &got); err != nil {
		return err
	}
	if err := expect.All(errs.Wrap(number.Equal(len(got.Readings), 1), "date must have one actual reading")); err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			number.Equal(got.Readings[0].DayKW, dayKW),
			"field day_kw must be the last sent value",
		),
		errs.Wrap(
			number.Equal(got.Readings[0].NightKW, nightKW),
			"field night_kw must be the last sent value",
		),
	)
}

package correction

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

// Run requests the readings of the date sent by Post and checks that exactly one row with the sent values is returned.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.Require(fixture.RealtyID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s/realty/%s/readings?date_from=%s&date_to=%s",
			s.baseURL,
			id,
			readingDate,
			readingDate,
		),
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
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	var got struct {
		Readings []struct {
			DayKW   uint32 `json:"day_kw"`
			NightKW uint32 `json:"night_kw"`
		} `json:"readings"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return err
	}
	err = expect.All(
		errs.Wrap(
			number.Equal(len(got.Readings), 1),
			"date must have one actual reading",
		),
	)
	if err != nil {
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

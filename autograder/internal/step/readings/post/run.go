package post

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

// Run sends a reading for the object stored under the realty id and checks the response.
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

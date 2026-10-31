// Package tariffs provides the steps that check the tariff endpoint.
package tariffs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// lookupDate is the date the tariff in force is requested for.
const lookupDate = "2026-10-27"

// Get checks that GET /tariffs?date= returns the tariff version in force on the date.
type Get struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewGet creates a Get step that sends its request to baseURL through cli.
func NewGet(cli httpx.Doer, baseURL string) Get {
	return Get{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Get) Name() string {
	return "tariffs get"
}

// Description returns a brief explanation of what the step checks.
func (s Get) Description() string {
	return "Checks that GET /tariffs?date= returns the tariff version in force on the date."
}

// Run requests the tariff in force on lookupDate and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Get) Run(ctx context.Context, _ step.Variables) error {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/tariffs?date=%s", s.baseURL, lookupDate),
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
		ID             uint64 `json:"id"`
		Name           string `json:"name"`
		PriceDayKWPH   uint32 `json:"price_day_kwph"`
		PriceNightKWPH uint32 `json:"price_night_kwph"`
		ValidFrom      string `json:"valid_from"`
		ValidTo        string `json:"valid_to"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.NotEmpty(got.Name),
			"field name",
		),
		errs.Wrap(
			number.Positive(got.PriceDayKWPH),
			"field price_day_kwph",
		),
		errs.Wrap(
			number.Positive(got.PriceNightKWPH),
			"field price_night_kwph",
		),
		errs.Wrap(
			strings.NotEmpty(got.ValidFrom),
			"field valid_from",
		),
		errs.Wrap(
			strings.NotEmpty(got.ValidTo),
			"field valid_to",
		),
	)
}

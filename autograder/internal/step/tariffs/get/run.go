package get

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

// Run requests the tariff in force on the lookup date and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	_ step.Variables,
) error {
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
		Name           string `json:"name"`
		ValidFrom      string `json:"valid_from"`
		ValidTo        string `json:"valid_to"`
		ID             uint64 `json:"id"`
		PriceDayKWPH   uint32 `json:"price_day_kwph"`
		PriceNightKWPH uint32 `json:"price_night_kwph"`
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

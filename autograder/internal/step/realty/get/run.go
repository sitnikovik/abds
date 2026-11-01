package get

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run requests the object stored under the realty id and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
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
		s.baseURL+"/realty/"+id,
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
		Address string `json:"address"`
		Type    string `json:"type"`
		ID      uint64 `json:"id"`
		Square  uint32 `json:"square"`
		Floor   uint8  `json:"floor"`
		Porch   uint8  `json:"porch"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(strconv.FormatUint(got.ID, 10), id),
			"field id",
		),
		errs.Wrap(
			strings.NotEmpty(got.Address),
			"field address",
		),
		errs.Wrap(
			number.Positive(got.Square),
			"field square",
		),
		errs.Wrap(
			number.Positive(got.Floor),
			"field floor",
		),
		errs.Wrap(
			number.Positive(got.Porch),
			"field porch",
		),
		errs.Wrap(
			expect.Any(
				strings.Equal(got.Type, "flat"),
				strings.Equal(got.Type, "nrcp"),
			),
			"field type must be flat or nrcp",
		),
	)
}

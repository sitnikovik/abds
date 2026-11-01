package badrequest

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

// Run requests the readings with a malformed date and checks the error response.
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
		s.baseURL+"/realty/"+id+"/readings?date_from=not-a-date&date_to=2026-10-31",
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
	var got struct {
		Message string `json:"message"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	return expect.All(
		httpx.Status(rsp, http.StatusBadRequest),
		httpx.JSONContent(rsp),
		errs.Wrap(
			strings.NotEmpty(got.Message),
			"field message",
		),
	)
}

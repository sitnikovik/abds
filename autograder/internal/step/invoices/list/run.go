package list

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/boolean"
	"github.com/sitnikovik/buggles/expect/number"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run requests the first page of invoices of the object and checks the response.
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
		fmt.Sprintf("%s/realty/%s/invoices?limit=%d&offset=0", s.baseURL, id, pageLimit),
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
		Invoices []struct {
			Created string `json:"created_at"`
		} `json:"invoices"`
		Count int `json:"count"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return err
	}
	sorted := true
	for i := 1; i < len(got.Invoices); i++ {
		if got.Invoices[i-1].Created < got.Invoices[i].Created {
			sorted = false
		}
	}
	return expect.All(
		errs.Wrap(
			number.LessOrEqual(len(got.Invoices), pageLimit),
			"page must not exceed the limit",
		),
		errs.Wrap(
			number.Positive(got.Count),
			"object must have invoices",
		),
		errs.Wrap(
			boolean.True(sorted),
			"invoices must be sorted by created_at descending",
		),
	)
}

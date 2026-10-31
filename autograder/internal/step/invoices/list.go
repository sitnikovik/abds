// Package invoices provides the steps that check the invoice endpoints.
package invoices

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

// pageLimit is the page size the step requests.
const pageLimit = 10

// List checks that GET /realty/{id}/invoices returns a page of invoices with the total count, newest first.
type List struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewList creates a List step that sends its request to baseURL through cli.
func NewList(cli httpx.Doer, baseURL string) List {
	return List{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s List) Name() string {
	return "invoices list"
}

// Description returns a brief explanation of what the step checks.
func (s List) Description() string {
	return "Checks that GET /realty/{id}/invoices returns a page of invoices, the total count and newest first order."
}

// Run requests the first page of invoices of the object and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s List) Run(ctx context.Context, vars step.Variables) error {
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
	if err := expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp)); err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	var got struct {
		Invoices []struct {
			Created string `json:"created_at"`
		} `json:"invoices"`
		Count int `json:"count"`
	}
	if err := httpx.Decode(body, &got); err != nil {
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

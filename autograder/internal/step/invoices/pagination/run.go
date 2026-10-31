package pagination

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
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run requests three pages of invoices and checks their sizes and that the second and the first pages differ.
//
// It returns an error if a variable is missing, a request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.Require(fixture.RealtyID)
	if err != nil {
		return err
	}
	defaults, err := s.Page(ctx, id, "")
	if err != nil {
		return err
	}
	first, err := s.Page(ctx, id, "?limit=1&offset=0")
	if err != nil {
		return err
	}
	second, err := s.Page(ctx, id, "?limit=1&offset=1")
	if err != nil {
		return err
	}
	if err := expect.All(
		errs.Wrap(
			number.LessOrEqual(len(defaults), defaultPageLimit),
			"default page must not exceed 10 invoices",
		),
		errs.Wrap(
			number.Equal(len(first), 1),
			"page of limit 1 must have one invoice",
		),
		errs.Wrap(
			number.Equal(len(second), 1),
			"second page of limit 1 must have one invoice",
		),
	); err != nil {
		return err
	}
	return errs.Wrap(
		strings.NotEqual(first[0], second[0]),
		"offset must move to another invoice",
	)
}

// Page requests the invoices of the object with the query string and returns their created_at values.
func (s Step) Page(ctx context.Context, id, query string) ([]string, error) {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/realty/%s/invoices%s", s.baseURL, id, query),
		nil,
	)
	if err != nil {
		return nil, err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return nil, err
	}
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return nil, errs.Wrap(err, "unexpected response")
	}
	var got struct {
		Invoices []struct {
			Created string `json:"created_at"`
		} `json:"invoices"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return nil, err
	}
	created := make([]string, 0, len(got.Invoices))
	for _, inv := range got.Invoices {
		created = append(created, inv.Created)
	}
	return created, nil
}

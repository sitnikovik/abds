package phonehistory

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

// Run changes the phone of the owner and checks that the previous phone is kept in phone_prev.
//
// It returns an error if a query fails, a request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.RequireUint64(fixture.UserID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf("SELECT phone FROM %s.dim_users FINAL WHERE id = ?", s.database),
		1,
		id,
	)
	if err != nil {
		return fmt.Errorf("read owner: %w", err)
	}
	err = expect.All(errs.Wrap(number.Positive(len(rows)), "owner must be stored"))
	if err != nil {
		return err
	}
	old := rows[0][0]
	next := "+79999999999"
	if old == next {
		next = "+79999999998"
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%d", s.baseURL, id),
		map[string]string{"phone": next},
	)
	if err != nil {
		return err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	patchBody, err := httpx.ReadBody(rsp)
	_ = rsp.Body.Close()
	if err != nil {
		return err
	}
	err = expect.All(
		httpx.Status(rsp, http.StatusNoContent),
		errs.Wrap(strings.Empty(string(patchBody)), "body must be empty"),
	)
	if err != nil {
		return err
	}
	got, err := s.Fetch(ctx, id)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(got.Phone, next),
			"field phone",
		),
		errs.Wrap(
			strings.Equal(got.PhonePrev, old),
			"field phone_prev must keep the previous phone",
		),
	)
}

// Fetch returns the owner from GET /users/{id}.
func (s Step) Fetch(ctx context.Context, id uint64) (Owner, error) {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/%d", s.baseURL, id),
		nil,
	)
	if err != nil {
		return Owner{}, err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return Owner{}, fmt.Errorf("send request: %w", err)
	}
	body, err := httpx.ReadBody(rsp)
	_ = rsp.Body.Close()
	if err != nil {
		return Owner{}, err
	}
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return Owner{}, errs.Wrap(err, "unexpected response")
	}
	var got Owner
	err = httpx.Decode(body, &got)
	if err != nil {
		return Owner{}, err
	}
	return got, nil
}

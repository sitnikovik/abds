// Package users provides the steps that check the user endpoints.
package users

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Get checks that GET /users/{id} returns the owner with the given id.
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
	return "users get"
}

// Description returns a brief explanation of what the step checks.
func (s Get) Description() string {
	return "Checks that GET /users/{id} returns the owner with the expected fields."
}

// Run requests the owner stored under the user id and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Get) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.Require(fixture.UserID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(ctx, http.MethodGet, s.baseURL+"/users/"+id, nil)
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
		ID    uint64 `json:"id"`
		FIO   string `json:"fio"`
		Phone string `json:"phone"`
	}
	if err := httpx.Decode(body, &got); err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(strconv.FormatUint(got.ID, 10), id),
			"field id",
		),
		errs.Wrap(
			strings.NotEmpty(got.FIO),
			"field fio",
		),
		errs.Wrap(
			strings.NotEmpty(got.Phone),
			"field phone",
		),
	)
}

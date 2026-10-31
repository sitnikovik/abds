package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// patchPhone is the phone the step sets on the owner.
const patchPhone = "+79999999999"

// Patch checks that PATCH /users/{id} changes the phone of the owner and answers 204 without a body.
type Patch struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewPatch creates a Patch step that sends its request to baseURL through cli.
func NewPatch(cli httpx.Doer, baseURL string) Patch {
	return Patch{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Patch) Name() string {
	return "users patch"
}

// Description returns a brief explanation of what the step checks.
func (s Patch) Description() string {
	return "Checks that PATCH /users/{id} updates the phone and answers 204 No Content."
}

// Run sends the new phone of the owner and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Patch) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.Require(fixture.UserID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPatch,
		s.baseURL+"/users/"+id,
		map[string]string{"phone": patchPhone},
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
		httpx.Status(rsp, http.StatusNoContent),
		strings.Empty(string(body)),
	)
}

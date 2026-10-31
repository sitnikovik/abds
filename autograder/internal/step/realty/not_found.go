package realty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// missingRealtyID is an id that no object has.
const missingRealtyID = "999999999"

// NotFound checks that GET /realty/{id} answers 404 with a message for an unknown object.
type NotFound struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewNotFound creates a NotFound step that sends its request to baseURL through cli.
func NewNotFound(cli httpx.Doer, baseURL string) NotFound {
	return NotFound{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s NotFound) Name() string {
	return "realty get not found"
}

// Description returns a brief explanation of what the step checks.
func (s NotFound) Description() string {
	return "Checks that GET /realty/{id} answers 404 with a message for an unknown object."
}

// Run requests an object that does not exist and checks the error response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s NotFound) Run(ctx context.Context, _ step.Variables) error {
	req, err := httpx.NewRequest(ctx, http.MethodGet, s.baseURL+"/realty/"+missingRealtyID, nil)
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
		httpx.Status(rsp, http.StatusNotFound),
		httpx.JSONContent(rsp),
		errs.Wrap(
			strings.NotEmpty(got.Message),
			"field message",
		),
	)
}

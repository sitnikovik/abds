package tariffs

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

// BadRequest checks that GET /tariffs without a date answers 400 with a message.
type BadRequest struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewBadRequest creates a BadRequest step that sends its request to baseURL through cli.
func NewBadRequest(cli httpx.Doer, baseURL string) BadRequest {
	return BadRequest{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s BadRequest) Name() string {
	return "tariffs bad request"
}

// Description returns a brief explanation of what the step checks.
func (s BadRequest) Description() string {
	return "Checks that GET /tariffs without a date answers 400 with a message."
}

// Run requests the tariffs without a date and checks the error response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s BadRequest) Run(ctx context.Context, _ step.Variables) error {
	req, err := httpx.NewRequest(ctx, http.MethodGet, s.baseURL+"/tariffs", nil)
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
	if err := httpx.Decode(body, &got); err != nil {
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

package notfound

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

// Run requests an object that does not exist and checks the error response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	_ step.Variables,
) error {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		s.baseURL+"/realty/"+missingRealtyID,
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
		httpx.Status(rsp, http.StatusNotFound),
		httpx.JSONContent(rsp),
		errs.Wrap(
			strings.NotEmpty(got.Message),
			"field message",
		),
	)
}

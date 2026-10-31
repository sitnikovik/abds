package ping

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Run sends GET to the ping endpoint and checks the status code and the body.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Step) Run(
	_ context.Context,
	_ step.Variables,
) error {
	rsp, err := s.client.Get(s.url)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	return expect.All(
		number.Equal(rsp.StatusCode, http.StatusOK),
		strings.Empty(string(body)),
	)
}

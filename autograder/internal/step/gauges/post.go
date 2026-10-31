// Package gauges provides the steps that check the gauge installation endpoint.
package gauges

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// installBody is the JSON body of POST /realty/{id}/gauges.
type installBody struct {
	// GaugeID is the id of the gauge to install on the object.
	GaugeID uint64 `json:"gauge_id"`
}

// Post checks that POST /realty/{id}/gauges installs a gauge and answers 201 without a body.
type Post struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewPost creates a Post step that sends its request to baseURL through cli.
func NewPost(cli httpx.Doer, baseURL string) Post {
	return Post{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Post) Name() string {
	return "gauges post"
}

// Description returns a brief explanation of what the step checks.
func (s Post) Description() string {
	return "Checks that POST /realty/{id}/gauges installs a gauge and answers 201 Created."
}

// Run installs the second gauge on the object and checks the response.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Post) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.Require(fixture.RealtyID)
	if err != nil {
		return err
	}
	gaugeID, err := vars.RequireUint64(fixture.NewGaugeID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/realty/%s/gauges", s.baseURL, id),
		installBody{GaugeID: gaugeID},
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
		httpx.Status(rsp, http.StatusCreated),
		errs.Wrap(
			strings.Empty(string(body)),
			"body must be empty",
		),
	)
}

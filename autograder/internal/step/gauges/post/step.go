// Package post provides the step that checks the gauge installation.
package post

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// installBody is the JSON body of POST /realty/{id}/gauges.
type installBody struct {
	// GaugeID is the id of the gauge to install on the object.
	GaugeID uint64 `json:"gauge_id"`
}

// Step checks that POST /realty/{id}/gauges installs a gauge and answers 201 without a body.
type Step struct {
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that sends its request to baseURL through cli.
func NewStep(cli httpx.Doer, baseURL string) Step {
	return Step{
		cli:     cli,
		baseURL: baseURL,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "gauges post"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that POST /realty/{id}/gauges installs a gauge and answers 201 Created."
}

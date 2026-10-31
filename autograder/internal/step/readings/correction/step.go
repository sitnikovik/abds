// Package correction provides the step that checks the reading correction.
package correction

import (
	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// readingDate is the date of the reading the step checks.
const readingDate = "2026-10-01"

// dayKW is the day counter value the step expects.
const dayKW = 1250

// nightKW is the night counter value the step expects.
const nightKW = 300

// Step checks that the readings of the date show only the last value sent for it (ReplacingMergeTree).
//
// It must run after Post, which sends the reading of the date.
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
	return "readings correction"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that the readings of a date show one actual value after a correction."
}

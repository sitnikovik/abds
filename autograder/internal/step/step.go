// Package step defines the checks of the autograder as ordered units of work.
package step

import (
	"context"
)

// Step is a single check of the autograder.
type Step interface {
	// Name returns the short name of the step.
	Name() string
	// Description returns a brief explanation of what the step checks.
	Description() string
	// Run executes the check and returns an error if the check fails.
	//
	// Values produced by the check are stored in vars for the following steps.
	Run(ctx context.Context, vars Variables) error
}

// Package step provides fake steps for tests.
package step

import (
	"context"

	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Step is a fake step that returns a predefined error and remembers whether it ran.
type Step struct {
	// err is returned by Run.
	err error
	// name is the name of the step.
	name string
	// ran reports whether Run was called.
	ran bool
}

// NewStep creates a fake step with the given name that returns err from Run.
func NewStep(name string, err error) *Step {
	return &Step{
		name: name,
		err:  err,
	}
}

// Name returns the name of the fake step.
func (s *Step) Name() string {
	return s.name
}

// Description returns a fixed description of the fake step.
func (s *Step) Description() string {
	return "fake step " + s.name
}

// Run records the call and returns the error the fake was created with.
func (s *Step) Run(_ context.Context, _ step.Variables) error {
	s.ran = true
	return s.err
}

// Ran reports whether Run was called.
func (s *Step) Ran() bool {
	return s.ran
}

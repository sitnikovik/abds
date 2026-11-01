package listnotfound

import (
	"context"

	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Run sends the request and checks the error response.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	return s.Step.Run(ctx, vars)
}

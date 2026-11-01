package replacement

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run reads the active gauges of the object and checks that only the new gauge is active.
//
// It returns an error if a variable is missing, the query fails or the result does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	realty, err := vars.RequireUint64(fixture.RealtyID)
	if err != nil {
		return err
	}
	want, err := vars.Require(fixture.NewGaugeID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf(
			"SELECT toString(gauge_id) FROM %s.fact_gauge_bindings FINAL WHERE realty_id = ?",
			s.database,
		),
		1,
		realty,
	)
	if err != nil {
		return fmt.Errorf("read bindings: %w", err)
	}
	if err := expect.All(
		errs.Wrap(
			number.Equal(len(rows), 1),
			"object must have exactly one active gauge",
		),
	); err != nil {
		return err
	}
	return errs.Wrap(
		strings.Equal(rows[0][0], want),
		"active gauge must be the installed one",
	)
}

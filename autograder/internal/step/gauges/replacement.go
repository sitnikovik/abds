package gauges

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

// clickhouseClient is the part of the ClickHouse client the replacement step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Replacement checks that after POST /realty/{id}/gauges the new gauge is the only one installed on the object.
//
// It reads the bindings with FINAL, so the cancelled installation does not count.
// It must run after Post, which installs the gauge.
type Replacement struct {
	// db reads the bindings.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
}

// NewReplacement creates a Replacement step that reads the bindings from database through db.
func NewReplacement(db clickhouseClient, database string) Replacement {
	return Replacement{
		db:       db,
		database: database,
	}
}

// Name returns the name of the step.
func (s Replacement) Name() string {
	return "gauges replacement"
}

// Description returns a brief explanation of what the step checks.
func (s Replacement) Description() string {
	return "Checks that the installed gauge is replaced and the previous installation is cancelled."
}

// Run reads the active gauges of the object and checks that only the new gauge is active.
//
// It returns an error if a variable is missing, the query fails or the result does not match the expectations.
func (s Replacement) Run(ctx context.Context, vars step.Variables) error {
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
		fmt.Sprintf("SELECT toString(gauge_id) FROM %s.fact_gauge_bindings FINAL WHERE realty_id = ?", s.database),
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

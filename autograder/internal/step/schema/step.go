// Package schema provides the step that checks the tables of the student's database.
package schema

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Step checks that every table of lab 1 exists with the expected engine and columns.
type Step struct {
	// db runs the queries against the system tables.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
}

// NewStep creates a Step that inspects database through db.
func NewStep(db clickhouseClient, database string) Step {
	return Step{
		db:       db,
		database: database,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "schema"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that the tables exist with the expected engines and columns."
}

// Run reads the engines and the columns of the database and checks them against the lab spec.
//
// It returns an error if a query fails or any table or column does not match the spec.
func (s Step) Run(ctx context.Context, _ step.Variables) error {
	engines, err := s.db.Strings(ctx, enginesQuery, 2, s.database)
	if err != nil {
		return fmt.Errorf("read table engines: %w", err)
	}
	columns, err := s.db.Strings(ctx, columnsQuery, 3, s.database)
	if err != nil {
		return fmt.Errorf("read columns: %w", err)
	}
	engineOf := make(map[string]string, len(engines))
	for _, row := range engines {
		engineOf[row[0]] = row[1]
	}
	typeOf := make(map[string]string, len(columns))
	for _, row := range columns {
		typeOf[row[0]+"."+row[1]] = row[2]
	}
	checks := make([]error, 0)
	for _, t := range tables {
		checks = append(
			checks,
			errs.Wrap(
				strings.Equal(engineOf[t.name], t.engine),
				"table "+t.name+" engine",
			),
		)
		for _, c := range t.columns {
			checks = append(
				checks,
				errs.Wrap(
					strings.Equal(typeOf[t.name+"."+c.name], c.typ),
					"column "+t.name+"."+c.name+" type",
				),
			)
		}
	}
	return expect.All(checks[0], checks[1:]...)
}

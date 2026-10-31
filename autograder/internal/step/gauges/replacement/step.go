// Package replacement provides the step that checks the gauge replacement.
package replacement

import (
	"context"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Step checks that after POST /realty/{id}/gauges the new gauge is the only one installed on the object.
//
// It reads the bindings with FINAL, so the cancelled installation does not count.
// It must run after Post, which installs the gauge.
type Step struct {
	// db reads the bindings.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
}

// NewStep creates a Step that reads the bindings from database through db.
func NewStep(db clickhouseClient, database string) Step {
	return Step{
		db:       db,
		database: database,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "gauges replacement"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that the installed gauge is replaced and the previous installation is cancelled."
}

// Package schema provides the step that checks the tables of the student's database.
package schema

import (
	"context"
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

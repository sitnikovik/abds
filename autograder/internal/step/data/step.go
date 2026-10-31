// Package data provides the step that checks the volume of the data loaded by make setup.
package data

import (
	"context"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Count runs the query and returns the value of its first row.
	Count(ctx context.Context, query string, args ...any) (uint64, error)
}

// Step checks that make setup loaded enough rows into the tables of lab 1.
type Step struct {
	// db runs the count queries.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
}

// NewStep creates a Step that counts rows in database through db.
func NewStep(db clickhouseClient, database string) Step {
	return Step{
		db:       db,
		database: database,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "data"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that make setup loaded the expected volume of data."
}

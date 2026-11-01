// Package fio provides the step that checks the owner's full name update.
package fio

import (
	"context"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
	// Count runs the query and returns the value of its first row.
	Count(ctx context.Context, query string, args ...any) (uint64, error)
}

// Step checks that PATCH /users/{id} overwrites the full name and the previous name is not stored anymore.
type Step struct {
	// db reads the stored owner.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
	// cli sends the requests.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that reads the owner from database through db and sends its requests through cli.
func NewStep(db clickhouseClient, database string, cli httpx.Doer, baseURL string) Step {
	return Step{
		db:       db,
		database: database,
		cli:      cli,
		baseURL:  baseURL,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "users fio"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that PATCH /users/{id} overwrites the full name and the previous name is not stored."
}

// Package phonehistory provides the step that checks the phone history of the owner.
package phonehistory

import (
	"context"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Step checks that PATCH /users/{id} keeps the previous phone in phone_prev (SCD 3).
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
	return "users phone history"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that PATCH /users/{id} keeps the previous phone in phone_prev."
}

// Owner is the part of the owner returned by GET /users/{id} the step checks.
type Owner struct {
	// Phone is the current phone of the owner.
	Phone string `json:"phone"`
	// PhonePrev is the phone the owner had before the current one.
	PhonePrev string `json:"phone_prev"`
}

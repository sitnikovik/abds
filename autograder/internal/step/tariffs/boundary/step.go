// Package boundary provides the step that checks the tariff version boundary.
package boundary

import (
	"context"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

// dateTimeLayout is the layout of the DateTime values in the API and in the database.
const dateTimeLayout = "2006-01-02 15:04:05"

// dayLayout is the layout of the dates in the requests.
const dayLayout = "2006-01-02"

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Step checks that the tariff version changes on the day its second version starts (SCD 2).
type Step struct {
	// db reads the versions of the tariff.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
	// cli sends the requests.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewStep creates a Step that reads the versions from database through db and sends its requests through cli.
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
	return "tariffs boundary"
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return "Checks that the tariff version in force changes on the day its second version starts."
}

// Version is the part of the tariff version returned by GET /tariffs the step checks.
type Version struct {
	// ValidFrom is the start of the version.
	ValidFrom string `json:"valid_from"`
	// ValidTo is the end of the version.
	ValidTo string `json:"valid_to"`
}

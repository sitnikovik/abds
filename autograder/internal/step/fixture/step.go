// Package fixture provides the step that reads ids of the loaded data for the endpoint checks.
package fixture

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Keys of the variables the step stores for the endpoint checks.
const (
	// RealtyID is the id of an object of the company.
	RealtyID = "realty_id"
	// UserID is the id of an owner.
	UserID = "user_id"
	// TariffID is the id of a tariff.
	TariffID = "tariff_id"
	// GaugeID is the id of a gauge that is not used by the endpoint checks.
	GaugeID = "gauge_id"
	// NewGaugeID is the id of another gauge that the endpoint checks install on an object.
	NewGaugeID = "new_gauge_id"
)

// clickhouseClient is the part of the ClickHouse client the step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// source is a table the step takes ids from, and the variables the ids are stored under.
type source struct {
	// table is the name of the table.
	table string
	// keys are the variables the first ids of the table are stored under, in order.
	keys []string
}

// sources lists the tables the step reads ids from.
var sources = []source{
	{table: "dim_realty", keys: []string{RealtyID}},
	{table: "dim_users", keys: []string{UserID}},
	{table: "dim_tariffs", keys: []string{TariffID}},
	{table: "dim_gauges", keys: []string{GaugeID, NewGaugeID}},
}

// Step stores the ids of the loaded data in the variables for the endpoint checks.
type Step struct {
	// db runs the queries.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
}

// NewStep creates a Step that reads ids from database through db.
func NewStep(db clickhouseClient, database string) Step {
	return Step{
		db:       db,
		database: database,
	}
}

// Name returns the name of the step.
func (s Step) Name() string {
	return "fixture"
}

// Description returns a brief explanation of what the step does.
func (s Step) Description() string {
	return "Reads the ids of the loaded data for the endpoint checks."
}

// Run reads the first ids of every source table and stores them in vars.
//
// It returns an error if a query fails or a table has fewer rows than its variables need.
func (s Step) Run(ctx context.Context, vars step.Variables) error {
	for _, src := range sources {
		query := fmt.Sprintf(
			"SELECT toString(id) FROM %s.%s ORDER BY id LIMIT %d",
			s.database,
			src.table,
			len(src.keys),
		)
		rows, err := s.db.Strings(ctx, query, 1)
		if err != nil {
			return fmt.Errorf("read ids of %s: %w", src.table, err)
		}
		err = expect.All(
			errs.Wrap(
				number.GreaterOrEqual(len(rows), len(src.keys)),
				fmt.Sprintf("rows in %s", src.table),
			),
		)
		if err != nil {
			return err
		}
		for i, key := range src.keys {
			vars.Set(key, rows[i][0])
		}
	}
	return nil
}

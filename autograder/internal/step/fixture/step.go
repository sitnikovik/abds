// Package fixture provides the step that reads ids of the loaded data for the endpoint checks.
package fixture

import (
	"context"
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

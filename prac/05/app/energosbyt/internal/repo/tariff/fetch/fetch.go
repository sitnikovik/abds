package fetch

import (
	chgo "github.com/ClickHouse/clickhouse-go/v2"

	"energosbyt/internal/repo/gauge/filter"
)

// Filter defines the interface for building WHERE clause filters.
type Filter interface {
	// Where returns the WHERE clause string.
	Where() string

	// Empty returns true if no filter conditions have been added.
	Empty() bool
}

// Fetch manages the retrieval of tariffs from ClickHouse with optional filtering and limiting.
type Fetch struct {
	// db stores the ClickHouse database connection.
	db chgo.Conn

	// ftr stores the filter for building WHERE clauses.
	ftr Filter

	// limit stores the maximum number of results to return.
	limit int
}

// NewFetch creates a new Fetch instance with the given database connection and options.
func NewFetch(
	db chgo.Conn,
	opts ...Option,
) Fetch {
	t := Fetch{
		db:  db,
		ftr: &filter.Filter{},
	}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

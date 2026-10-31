// Package clickhouse provides a read-only client over the ClickHouse native protocol.
package clickhouse

import (
	"context"
	"errors"
	"fmt"

	chdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// conn is the part of the ClickHouse connection the Client uses.
type conn interface {
	// Query runs the query with args and returns the rows it selects.
	Query(ctx context.Context, query string, args ...any) (chdriver.Rows, error)
}

// Client runs queries against a ClickHouse server.
type Client struct {
	// conn runs the queries.
	conn conn
}

// NewClient creates a Client that runs its queries through conn.
func NewClient(conn conn) Client {
	return Client{
		conn: conn,
	}
}

// Count runs the query and returns the value of its first row.
//
// The query must select exactly one UInt64 column, for example "SELECT count() FROM t".
// It returns an error if the query fails or returns no rows.
func (c Client) Count(ctx context.Context, query string, args ...any) (uint64, error) {
	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("run query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return 0, fmt.Errorf("iterate rows: %w", err)
		}
		return 0, errors.New("query returned no rows")
	}
	var n uint64
	if err := rows.Scan(&n); err != nil {
		return 0, fmt.Errorf("scan count: %w", err)
	}
	return n, nil
}

// Strings runs the query and returns all rows as strings.
//
// Every row must have exactly columns string columns.
// It returns an error if the query fails or a row cannot be scanned.
func (c Client) Strings(
	ctx context.Context,
	query string,
	columns int,
	args ...any,
) ([][]string, error) {
	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("run query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	var out [][]string
	for rows.Next() {
		vals := make([]string, columns)
		dest := make([]any, columns)
		for i := range vals {
			dest[i] = &vals[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		out = append(out, vals)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}
	return out, nil
}

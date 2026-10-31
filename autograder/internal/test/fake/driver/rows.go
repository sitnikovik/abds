// Package driver provides fake ClickHouse driver primitives for tests.
package driver

import (
	"context"
	"errors"
	"fmt"

	chdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Rows is a fake driver.Rows that iterates over predefined values.
type Rows struct {
	// err is returned by Err after the iteration.
	err error
	// data holds the rows, each row holds one value per column.
	data [][]any
	// idx is the index of the current row.
	idx int
}

// NewRows creates a Rows that iterates over data and returns err from Err.
func NewRows(data [][]any, err error) *Rows {
	return &Rows{
		data: data,
		idx:  -1,
		err:  err,
	}
}

// Next advances to the next row and reports whether there is one.
func (r *Rows) Next() bool {
	r.idx++
	return r.idx < len(r.data)
}

// Scan copies the values of the current row into dest, which must be
// pointers to uint64 or string values matching the stored types.
func (r *Rows) Scan(dest ...any) error {
	if r.idx < 0 || r.idx >= len(r.data) {
		return errors.New("scan without a current row")
	}
	row := r.data[r.idx]
	if len(dest) != len(row) {
		return fmt.Errorf("scan %d values into %d destinations", len(row), len(dest))
	}
	for i, d := range dest {
		switch p := d.(type) {
		case *uint64:
			v, ok := row[i].(uint64)
			if !ok {
				return fmt.Errorf("value %d is %T, not uint64", i, row[i])
			}
			*p = v
		case *string:
			v, ok := row[i].(string)
			if !ok {
				return fmt.Errorf("value %d is %T, not string", i, row[i])
			}
			*p = v
		default:
			return fmt.Errorf("unsupported destination %T", d)
		}
	}
	return nil
}

// ScanStruct is not supported by the fake and always returns an error.
func (r *Rows) ScanStruct(_ any) error {
	return errors.New("scan struct is not supported")
}

// ColumnTypes returns nil because the fake does not describe its columns.
func (r *Rows) ColumnTypes() []chdriver.ColumnType {
	return nil
}

// Totals is not supported by the fake and always returns an error.
func (r *Rows) Totals(_ ...any) error {
	return errors.New("totals are not supported")
}

// Columns returns nil because the fake does not name its columns.
func (r *Rows) Columns() []string {
	return nil
}

// Close releases the rows and always succeeds.
func (r *Rows) Close() error {
	return nil
}

// Err returns the error the fake was created with.
func (r *Rows) Err() error {
	return r.err
}

// HasData reports whether there are rows in the fake.
func (r *Rows) HasData() bool {
	return len(r.data) > 0
}

// Ensure Rows satisfies the driver interface.
var _ chdriver.Rows = (*Rows)(nil)

// Conn is a fake connection that returns predefined rows for every query.
type Conn struct {
	// rows are returned by Query.
	rows *Rows
	// err is returned by Query instead of rows.
	err error
}

// NewConn creates a Conn whose Query returns rows, or err when err is not nil.
func NewConn(rows *Rows, err error) *Conn {
	return &Conn{
		rows: rows,
		err:  err,
	}
}

// Query returns the predefined rows, or the error the fake was created with.
func (c *Conn) Query(_ context.Context, _ string, _ ...any) (chdriver.Rows, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.rows, nil
}

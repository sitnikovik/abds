// Package clickhouse provides a fake ClickHouse client for tests of the steps.
//
// Each rule matches a query by a substring, and the first matching rule answers the query.
// A query without a matching rule fails with an error.
package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// countRule answers queries containing substr with n.
type countRule struct {
	// substr is the text the query must contain.
	substr string
	// n is the count returned for the query.
	n uint64
}

// stringsRule answers queries containing substr with rows.
type stringsRule struct {
	// substr is the text the query must contain.
	substr string
	// rows are the rows returned for the query.
	rows [][]string
}

// Client is a fake ClickHouse client that answers queries by matching rules.
type Client struct {
	// err is returned for every query when it is not nil.
	err error
	// counts are the rules for Count.
	counts []countRule
	// rows are the rules for Strings.
	rows []stringsRule
}

// Option configures a fake Client.
type Option func(*Client)

// WithCount answers Count for queries containing substr with n.
func WithCount(substr string, n uint64) Option {
	return func(c *Client) {
		c.counts = append(c.counts, countRule{substr: substr, n: n})
	}
}

// WithStrings answers Strings for queries containing substr with rows.
func WithStrings(substr string, rows [][]string) Option {
	return func(c *Client) {
		c.rows = append(c.rows, stringsRule{substr: substr, rows: rows})
	}
}

// WithError makes every query fail with err.
func WithError(err error) Option {
	return func(c *Client) {
		c.err = err
	}
}

// NewClient creates a fake Client configured by opts.
func NewClient(opts ...Option) *Client {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Count returns the count of the first Count rule matching query.
func (c *Client) Count(_ context.Context, query string, _ ...any) (uint64, error) {
	if c.err != nil {
		return 0, c.err
	}
	for _, r := range c.counts {
		if strings.Contains(query, r.substr) {
			return r.n, nil
		}
	}
	return 0, fmt.Errorf("no count rule matches query %q", query)
}

// Strings returns the rows of the first Strings rule matching query.
func (c *Client) Strings(_ context.Context, query string, _ int, _ ...any) ([][]string, error) {
	if c.err != nil {
		return nil, c.err
	}
	for _, r := range c.rows {
		if strings.Contains(query, r.substr) {
			return r.rows, nil
		}
	}
	return nil, errors.New("no strings rule matches query " + query)
}

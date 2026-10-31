package clickhouse_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	errs "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/driver"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/env"
)

func TestClient_Count(t *testing.T) {
	t.Parallel()

	t.Run("returns the value of the first row", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(
				driver.NewRows(
					[][]any{
						{uint64(42)},
					},
					nil,
				),
				nil,
			),
		)
		n, err := cli.Count(context.Background(), "SELECT count() FROM t")
		assert.NoError(
			t,
			expect.All(
				errs.NoError(err),
				number.Equal(n, uint64(42)),
			),
		)
	})
	t.Run("fails when the query fails", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(nil, errors.New("connection refused")),
		)
		_, err := cli.Count(context.Background(), "SELECT count() FROM t")
		assert.NoError(
			t,
			errs.Contains(err, "connection refused"),
		)
	})
	t.Run("fails when the query returns no rows", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(driver.NewRows(nil, nil), nil),
		)
		_, err := cli.Count(context.Background(), "SELECT count() FROM t")
		assert.NoError(
			t,
			errs.Contains(err, "no rows"),
		)
	})
	t.Run("fails when iteration fails", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(driver.NewRows(nil, errors.New("broken stream")), nil),
		)
		_, err := cli.Count(context.Background(), "SELECT count() FROM t")
		assert.NoError(
			t,
			errs.Contains(err, "broken stream"),
		)
	})
	t.Run("fails when the value cannot be scanned", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(
				driver.NewRows(
					[][]any{
						{"not a number"},
					},
					nil,
				),
				nil,
			),
		)
		_, err := cli.Count(context.Background(), "SELECT count() FROM t")
		assert.NoError(
			t,
			errs.Contains(err, "scan count"),
		)
	})
}

func TestClient_Strings(t *testing.T) {
	t.Parallel()

	t.Run("returns all rows", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(
				driver.NewRows(
					[][]any{
						{"dim_realty", "MergeTree"},
						{"dim_users", "ReplacingMergeTree"},
					},
					nil,
				),
				nil,
			),
		)
		rows, err := cli.Strings(context.Background(), "SELECT name, engine FROM system.tables", 2)
		assert.NoError(
			t,
			expect.All(
				errs.NoError(err),
				number.Equal(len(rows), 2),
				strings.Equal(rows[0][0], "dim_realty"),
				strings.Equal(rows[1][1], "ReplacingMergeTree"),
			),
		)
	})
	t.Run("returns no rows for an empty result", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(driver.NewRows(nil, nil), nil),
		)
		rows, err := cli.Strings(context.Background(), "SELECT name FROM system.tables", 1)
		assert.NoError(
			t,
			expect.All(
				errs.NoError(err),
				number.Equal(len(rows), 0),
			),
		)
	})
	t.Run("fails when the query fails", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(nil, errors.New("unknown table")),
		)
		_, err := cli.Strings(context.Background(), "SELECT name FROM missing", 1)
		assert.NoError(
			t,
			errs.Contains(err, "unknown table"),
		)
	})
	t.Run("fails when a row cannot be scanned", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(
				driver.NewRows(
					[][]any{
						{uint64(1)},
					},
					nil,
				),
				nil,
			),
		)
		_, err := cli.Strings(context.Background(), "SELECT id FROM t", 1)
		assert.NoError(
			t,
			errs.Contains(err, "scan row"),
		)
	})
	t.Run("fails when iteration fails", func(t *testing.T) {
		t.Parallel()
		cli := clickhouse.NewClient(
			driver.NewConn(driver.NewRows(nil, errors.New("broken stream")), nil),
		)
		_, err := cli.Strings(context.Background(), "SELECT name FROM system.tables", 1)
		assert.NoError(
			t,
			errs.Contains(err, "broken stream"),
		)
	})
}

func TestOpen(t *testing.T) {
	t.Parallel()

	t.Run("creates a client without contacting the server", func(t *testing.T) {
		t.Parallel()
		cfg, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "8080",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "1",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(t, errs.NoError(err))
		_, err = clickhouse.Open(cfg)
		assert.NoError(
			t,
			errs.NoError(err),
		)
	})
}

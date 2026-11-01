package data_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	impl "github.com/sitnikovik/abds/autograder/internal/step/data"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
)

// loadedCount is a substring of a count query and the count the fake returns for it.
type loadedCount struct {
	// substr is the text the query must contain.
	substr string
	// n is the count the fake returns.
	n uint64
}

// loadedCounts returns the counts of a database loaded with the minimum volumes.
//
// The rule for replaced bindings goes before the rule for all bindings, because both queries name the table.
var loadedCounts = []loadedCount{
	{substr: "dim_realty", n: 1000},
	{substr: "dim_gauges", n: 1000},
	{substr: "dim_users", n: 1000},
	{substr: "fact_meter_readings", n: 1000000},
	{substr: "fact_invoice", n: 36000},
	{substr: "version >= 2", n: 1},
	{substr: "fact_gauge_bindings", n: 1000},
	{substr: "dim_tariffs", n: 1},
}

// ClientWithout returns a fake database that answers every count except the ones containing skip.
func ClientWithout(skip string) *fakeclickhouse.Client {
	opts := make(
		[]fakeclickhouse.Option,
		0,
		len(loadedCounts),
	)
	for _, c := range loadedCounts {
		if c.substr != skip {
			opts = append(opts, fakeclickhouse.WithCount(c.substr, c.n))
		}
	}
	return fakeclickhouse.NewClient(opts...)
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when every table has enough rows", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout(""),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when readings are fewer than required", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithCount("dim_realty", 1000),
			fakeclickhouse.WithCount("dim_gauges", 1000),
			fakeclickhouse.WithCount("dim_users", 1000),
			fakeclickhouse.WithCount("fact_meter_readings", 999999),
			fakeclickhouse.WithCount("fact_invoice", 36000),
			fakeclickhouse.WithCount("version >= 2", 1),
			fakeclickhouse.WithCount("fact_gauge_bindings", 1000),
			fakeclickhouse.WithCount("dim_tariffs", 1),
		)
		err := impl.
			NewStep(
				db,
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "rows in fact_meter_readings"))
	})

	t.Run("fails when no gauge was replaced", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithCount("dim_realty", 1000),
			fakeclickhouse.WithCount("dim_gauges", 1000),
			fakeclickhouse.WithCount("dim_users", 1000),
			fakeclickhouse.WithCount("fact_meter_readings", 1000000),
			fakeclickhouse.WithCount("fact_invoice", 36000),
			fakeclickhouse.WithCount("version >= 2", 0),
			fakeclickhouse.WithCount("fact_gauge_bindings", 1000),
			fakeclickhouse.WithCount("dim_tariffs", 1),
		)
		err := impl.
			NewStep(
				db,
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "gauge bindings after a replacement"))
	})

	t.Run("fails when the tariffs have no versions", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithCount("dim_realty", 1000),
			fakeclickhouse.WithCount("dim_gauges", 1000),
			fakeclickhouse.WithCount("dim_users", 1000),
			fakeclickhouse.WithCount("fact_meter_readings", 1000000),
			fakeclickhouse.WithCount("fact_invoice", 36000),
			fakeclickhouse.WithCount("version >= 2", 1),
			fakeclickhouse.WithCount("fact_gauge_bindings", 1000),
			fakeclickhouse.WithCount("dim_tariffs", 0),
		)
		err := impl.
			NewStep(
				db,
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "tariffs with at least two versions"))
	})

	t.Run("fails when gauges are fewer than required", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithCount("dim_realty", 1000),
			fakeclickhouse.WithCount("dim_gauges", 999),
			fakeclickhouse.WithCount("dim_users", 1000),
			fakeclickhouse.WithCount("fact_meter_readings", 1000000),
			fakeclickhouse.WithCount("fact_invoice", 36000),
			fakeclickhouse.WithCount("version >= 2", 1),
			fakeclickhouse.WithCount("fact_gauge_bindings", 1000),
			fakeclickhouse.WithCount("dim_tariffs", 1),
		)
		err := impl.
			NewStep(
				db,
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "rows in dim_gauges"))
	})

	t.Run("fails when the realty count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("dim_realty"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count dim_realty"))
	})

	t.Run("fails when the gauges count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("dim_gauges"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count dim_gauges"))
	})

	t.Run("fails when the users count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("dim_users"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count dim_users"))
	})

	t.Run("fails when the bindings count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("fact_gauge_bindings"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count fact_gauge_bindings"))
	})

	t.Run("fails when the readings count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("fact_meter_readings"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count fact_meter_readings"))
	})

	t.Run("fails when the invoices count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("fact_invoice"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count fact_invoice"))
	})

	t.Run("fails when the tariffs count fails", func(t *testing.T) {
		t.Parallel()
		err := impl.
			NewStep(
				ClientWithout("dim_tariffs"),
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "count tariffs with versions"))
	})

	t.Run("fails when every count query fails", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(fakeclickhouse.WithError(errors.New("connection refused")))
		err := impl.
			NewStep(
				db,
				"energosbyt",
			).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(
			t,
			expect.All(
				checkerr.Error(err),
				checkerr.Contains(err, "count dim_realty"),
			),
		)
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := impl.NewStep(
			fakeclickhouse.NewClient(),
			"energosbyt",
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "data"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

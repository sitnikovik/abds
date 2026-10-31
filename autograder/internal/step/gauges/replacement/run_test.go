package replacement_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"
	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	impl "github.com/sitnikovik/abds/autograder/internal/step/gauges/replacement"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
)

// VarsWithRealtyAndNewGauge returns variables with the realty id 1 and the new gauge id 2.
func VarsWithRealtyAndNewGauge() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.NewGaugeID, "2")
	return vars
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when only the new gauge is active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"2"}}),
		)
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						db,
						"energosbyt",
					).
					Run(
						context.Background(),
						VarsWithRealtyAndNewGauge(),
					),
			),
		)
	})

	t.Run("fails when the old gauge is still active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"1"}, {"2"}}),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						db,
						"energosbyt",
					).
					Run(
						context.Background(),
						VarsWithRealtyAndNewGauge(),
					),
				"exactly one active gauge",
			),
		)
	})

	t.Run("fails when another gauge is active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"1"}}),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						db,
						"energosbyt",
					).
					Run(
						context.Background(),
						VarsWithRealtyAndNewGauge(),
					),
				"installed one",
			),
		)
	})

	t.Run("fails when the query fails", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(fakeclickhouse.WithError(errors.New("connection refused")))
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						db,
						"energosbyt",
					).
					Run(
						context.Background(),
						VarsWithRealtyAndNewGauge(),
					),
				"read bindings",
			),
		)
	})

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakeclickhouse.NewClient(),
						"energosbyt",
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				fixture.RealtyID,
			),
		)
	})

	t.Run("fails when the new gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakeclickhouse.NewClient(),
						"energosbyt",
					).
					Run(
						context.Background(),
						vars,
					),
				fixture.NewGaugeID,
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
				strings.Equal(s.Name(), "gauges replacement"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

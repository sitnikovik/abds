package gauges_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/gauges"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
)

// varsWithRealtyAndNewGauge returns variables with the realty id 1 and the new gauge id 2.
func varsWithRealtyAndNewGauge() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.NewGaugeID, "2")
	return vars
}

func TestReplacement_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when only the new gauge is active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"2"}}),
		)
		s := gauges.NewReplacement(db, "energosbyt")
		err := s.Run(context.Background(), varsWithRealtyAndNewGauge())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the old gauge is still active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"1"}, {"2"}}),
		)
		s := gauges.NewReplacement(db, "energosbyt")
		err := s.Run(context.Background(), varsWithRealtyAndNewGauge())
		assert.NoError(t, checkerr.Contains(err, "exactly one active gauge"))
	})
	t.Run("fails when another gauge is active", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("fact_gauge_bindings", [][]string{{"1"}}),
		)
		s := gauges.NewReplacement(db, "energosbyt")
		err := s.Run(context.Background(), varsWithRealtyAndNewGauge())
		assert.NoError(t, checkerr.Contains(err, "installed one"))
	})
	t.Run("fails when the query fails", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(fakeclickhouse.WithError(errors.New("connection refused")))
		s := gauges.NewReplacement(db, "energosbyt")
		err := s.Run(context.Background(), varsWithRealtyAndNewGauge())
		assert.NoError(t, checkerr.Contains(err, "read bindings"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewReplacement(fakeclickhouse.NewClient(), "energosbyt")
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("fails when the new gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		s := gauges.NewReplacement(fakeclickhouse.NewClient(), "energosbyt")
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.NewGaugeID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewReplacement(fakeclickhouse.NewClient(), "energosbyt")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "gauges replacement"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

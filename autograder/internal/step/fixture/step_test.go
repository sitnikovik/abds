package fixture_test

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
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("stores the first ids of every table", func(t *testing.T) {
		t.Parallel()
		s := fixture.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("dim_realty", [][]string{{"1"}}),
				fakeclickhouse.WithStrings("dim_users", [][]string{{"2"}}),
				fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"3"}}),
				fakeclickhouse.WithStrings("dim_gauges", [][]string{{"4"}, {"5"}}),
			),
			"energosbyt",
		)
		vars := step.NewVariables()
		err := s.Run(context.Background(), vars)
		realty, _ := vars.Get(fixture.RealtyID)
		gauge, _ := vars.Get(fixture.GaugeID)
		newGauge, _ := vars.Get(fixture.NewGaugeID)
		assert.NoError(
			t,
			expect.All(
				checkerr.NoError(err),
				strings.Equal(realty, "1"),
				strings.Equal(gauge, "4"),
				strings.Equal(newGauge, "5"),
			),
		)
	})
	t.Run("fails when a table has too few rows", func(t *testing.T) {
		t.Parallel()
		s := fixture.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("dim_realty", [][]string{{"1"}}),
				fakeclickhouse.WithStrings("dim_users", [][]string{{"2"}}),
				fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"3"}}),
				fakeclickhouse.WithStrings("dim_gauges", [][]string{{"4"}}),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.Contains(err, "rows in dim_gauges"),
		)
	})
	t.Run("fails when the query fails", func(t *testing.T) {
		t.Parallel()
		s := fixture.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithError(errors.New("connection refused")),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.Contains(err, "read ids of dim_realty"),
		)
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := fixture.NewStep(fakeclickhouse.NewClient(), "energosbyt")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "fixture"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

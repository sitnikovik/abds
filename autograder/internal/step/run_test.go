package step_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/boolean"
	errs "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	fakestep "github.com/sitnikovik/abds/autograder/internal/test/fake/step"
)

func TestVariables(t *testing.T) {
	t.Parallel()

	t.Run("returns the stored value", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set("realty_id", "1")
		val, ok := vars.Lookup("realty_id")
		assert.NoError(
			t,
			expect.All(
				boolean.True(ok),
				strings.Equal(val, "1"),
			),
		)
	})
	t.Run("reports a missing key", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		_, ok := vars.Lookup("missing")
		assert.NoError(
			t,
			boolean.False(ok),
		)
	})
	t.Run("requires a stored value", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set("realty_id", "7")
		val, err := vars.Require("realty_id")
		assert.NoError(
			t,
			expect.All(
				errs.NoError(err),
				strings.Equal(val, "7"),
			),
		)
	})
	t.Run("fails to require a missing value", func(t *testing.T) {
		t.Parallel()
		_, err := step.NewVariables().Require("missing")
		assert.NoError(
			t,
			errs.Contains(err, "missing"),
		)
	})
	t.Run("requires an unsigned integer", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set("user_id", "42")
		n, err := vars.RequireUint64("user_id")
		assert.NoError(
			t,
			expect.All(
				errs.NoError(err),
				number.Equal(n, uint64(42)),
			),
		)
	})
	t.Run("fails to require a value that is not an integer", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set("user_id", "abc")
		_, err := vars.RequireUint64("user_id")
		assert.NoError(
			t,
			errs.Contains(err, "not an unsigned integer"),
		)
	})
	t.Run("fails to require an unsigned integer that is missing", func(t *testing.T) {
		t.Parallel()
		_, err := step.NewVariables().RequireUint64("user_id")
		assert.NoError(
			t,
			errs.Contains(err, "user_id"),
		)
	})
}

func TestRunAll(t *testing.T) {
	t.Parallel()

	t.Run("runs every step and reports the results in order", func(t *testing.T) {
		t.Parallel()
		steps := []step.Step{
			fakestep.NewStep("first", nil),
			fakestep.NewStep("second", errors.New("broken")),
			fakestep.NewStep("third", nil),
		}
		results := step.RunAll(context.Background(), steps)
		assert.NoError(
			t,
			expect.All(
				number.Equal(len(results), 3),
				strings.Equal(results[0].Name(), "first"),
				strings.Equal(results[1].Name(), "second"),
				strings.Equal(results[2].Name(), "third"),
				errs.NoError(results[0].Err()),
				errs.Error(results[1].Err()),
				errs.NoError(results[2].Err()),
			),
		)
	})
	t.Run("does not stop after a failed step", func(t *testing.T) {
		t.Parallel()
		last := fakestep.NewStep("last", nil)
		step.RunAll(
			context.Background(),
			[]step.Step{
				fakestep.NewStep("failing", errors.New("broken")),
				last,
			},
		)
		assert.NoError(
			t,
			boolean.True(last.Ran()),
		)
	})
	t.Run("returns no results for no steps", func(t *testing.T) {
		t.Parallel()
		results := step.RunAll(context.Background(), nil)
		assert.NoError(
			t,
			number.Equal(len(results), 0),
		)
	})
}

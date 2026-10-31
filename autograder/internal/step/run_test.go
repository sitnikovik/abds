package step_test

import (
	"testing"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/boolean"
	errs "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"
	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/abds/autograder/internal/step"
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

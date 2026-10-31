package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/errs"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("returns nil for a nil error", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			checkerr.NoError(errs.Wrap(nil, "table dim_realty")),
		)
	})
	t.Run("annotates the error with the message", func(t *testing.T) {
		t.Parallel()
		err := errs.Wrap(errors.New("boom"), "table dim_realty")
		assert.NoError(
			t,
			expect.All(
				checkerr.Error(err),
				checkerr.Contains(err, "table dim_realty"),
				checkerr.Contains(err, "boom"),
			),
		)
	})
}

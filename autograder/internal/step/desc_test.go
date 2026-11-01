package step_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
)

func TestDesc(t *testing.T) {
	t.Parallel()

	t.Run("keeps the name and the description", func(t *testing.T) {
		t.Parallel()
		d := step.NewDesc("users get", "Checks the owner.")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(d.Name(), "users get"),
				strings.Equal(d.Description(), "Checks the owner."),
			),
		)
	})
}

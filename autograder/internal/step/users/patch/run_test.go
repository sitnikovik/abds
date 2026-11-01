package patch_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"
	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	impl "github.com/sitnikovik/abds/autograder/internal/step/users/patch"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithUser returns variables that hold the user id 1.
func VarsWithUser() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.UserID, "1")
	return vars
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 204 with an empty body", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusNoContent,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
			),
		)
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"status code must be 204",
			),
		)
	})

	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Error(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusNoContent,
							"x",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
			),
		)
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							0,
							"",
							errors.New("connection refused"),
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"connection refused",
			),
		)
	})

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusNoContent,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				fixture.UserID,
			),
		)
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		var noerr error
		s := impl.NewStep(
			fakehttp.NewClient(
				http.StatusNoContent,
				"",
				noerr,
			),
			baseURL,
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users patch"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

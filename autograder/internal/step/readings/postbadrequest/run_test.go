package postbadrequest_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/readings/postbadrequest"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithRealty returns variables that hold the realty id 1.
func VarsWithRealty() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	return vars
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 400 with a message", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusBadRequest,
							`{"message":"invalid reading_date"}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
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
							`{"message":"x"}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"status code must be 400",
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
						VarsWithRealty(),
					),
				"connection refused",
			),
		)
	})

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusBadRequest,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				fixture.RealtyID,
			),
		)
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		var noerr error
		s := impl.NewStep(
			fakehttp.NewClient(
				http.StatusBadRequest,
				"",
				noerr,
			),
			baseURL,
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings post bad request"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

package correction_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/readings/correction"
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

	t.Run("passes on one reading with the last sent values", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"readings":[{"day_kw":1250,"night_kw":300}],"count":1}`,
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

	t.Run("fails when the date has two readings", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"readings":[{"day_kw":1250,"night_kw":300},{"day_kw":1,"night_kw":1}],"count":2}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"one actual reading",
			),
		)
	})

	t.Run("fails when the day value is an old one", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"readings":[{"day_kw":1,"night_kw":300}],"count":1}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"field day_kw",
			),
		)
	})

	t.Run("fails when the night value is an old one", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"readings":[{"day_kw":1250,"night_kw":1}],"count":1}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"field night_kw",
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
							http.StatusNotFound,
							`{"message":"x"}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"unexpected response",
			),
		)
	})

	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"decode body",
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
							http.StatusOK,
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
				http.StatusOK,
				"",
				noerr,
			),
			baseURL,
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings correction"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

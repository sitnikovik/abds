package post_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/gauges/post"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithRealtyAndGauge returns variables with a realty id and a new gauge id.
func VarsWithRealtyAndGauge() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.NewGaugeID, "2")
	return vars
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 201 with an empty body", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusCreated,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealtyAndGauge(),
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
							http.StatusBadRequest,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealtyAndGauge(),
					),
				"status code must be 201",
			),
		)
	})

	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusCreated,
							"{}",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealtyAndGauge(),
					),
				"body must be empty",
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
						VarsWithRealtyAndGauge(),
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
							http.StatusCreated,
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

	t.Run("fails when the gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusCreated,
							"",
							noerr,
						),
						baseURL,
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
		var noerr error
		s := impl.NewStep(
			fakehttp.NewClient(
				http.StatusCreated,
				"",
				noerr,
			),
			baseURL,
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "gauges post"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/readings/post"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithAll returns variables with every id the step needs.
func VarsWithAll() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.UserID, "1")
	vars.Set(fixture.GaugeID, "1")
	vars.Set(fixture.TariffID, "1")
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
						VarsWithAll(),
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
						VarsWithAll(),
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
						VarsWithAll(),
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
						VarsWithAll(),
					),
				"connection refused",
			),
		)
	})

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.GaugeID, "1")
		vars.Set(fixture.TariffID, "1")
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
				fixture.UserID,
			),
		)
	})

	t.Run("fails when the gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.UserID, "1")
		vars.Set(fixture.TariffID, "1")
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
				fixture.GaugeID,
			),
		)
	})

	t.Run("fails when the tariff id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.UserID, "1")
		vars.Set(fixture.GaugeID, "1")
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
				fixture.TariffID,
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
				strings.Equal(s.Name(), "readings post"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

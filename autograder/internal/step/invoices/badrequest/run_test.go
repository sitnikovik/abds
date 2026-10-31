package badrequest_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/invoices/badrequest"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 400 with a message", func(t *testing.T) {
		t.Parallel()
		var noerr error
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusBadRequest,
							`{"message":"limit must be positive"}`,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
			),
		)
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		var noerr error
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
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
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"status code must be 400",
			),
		)
	})

	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		var noerr error
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusBadRequest,
							`x`,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"unexpected response",
			),
		)
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
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
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
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
						"http://localhost:8080",
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
			"http://localhost:8080",
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "invoices bad request"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

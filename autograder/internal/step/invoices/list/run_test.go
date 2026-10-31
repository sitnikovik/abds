package list_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/invoices/list"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on a page with the total count", func(t *testing.T) {
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
							http.StatusOK,
							`{"invoices":[{"price":93000}],"count":36}`,
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

	t.Run("passes when invoices are sorted newest first", func(t *testing.T) {
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
							http.StatusOK,
							`{"invoices":[{"created_at":"2026-10-27 10:00:00"},{"created_at":"2026-09-27 10:00:00"}],"count":2}`,
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

	t.Run("fails when the page exceeds the limit", func(t *testing.T) {
		t.Parallel()
		var noerr error
		body := `{"invoices":[{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1}],"count":36}`
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							body,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"page must not exceed the limit",
			),
		)
	})

	t.Run("fails when the object has no invoices", func(t *testing.T) {
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
							`{"invoices":[],"count":0}`,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"object must have invoices",
			),
		)
	})

	t.Run("fails when invoices are sorted oldest first", func(t *testing.T) {
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
							`{"invoices":[{"created_at":"2026-09-27 10:00:00"},{"created_at":"2026-10-27 10:00:00"}],"count":2}`,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"sorted by created_at",
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
							http.StatusNotFound,
							`{"message":"x"}`,
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

	t.Run("fails on a body that is not JSON", func(t *testing.T) {
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
							`{`,
							noerr,
						),
						"http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
				"decode body",
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
							http.StatusOK,
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
				http.StatusOK,
				"",
				noerr,
			),
			"http://localhost:8080",
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "invoices list"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

package list_test

import (
	"context"
	"net/http"
	"testing"

	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	impl "github.com/sitnikovik/abds/autograder/internal/step/invoices/list"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

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
						), "http://localhost:8080",
					).
					Run(
						context.Background(),
						vars,
					),
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
}

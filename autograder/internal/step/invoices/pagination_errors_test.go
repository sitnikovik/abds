package invoices_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestPagination_RunErrors(t *testing.T) {
	t.Parallel()

	t.Run("fails when a page is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusOK,
			`{`,
			nil,
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})

	t.Run("fails when a page request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(fakehttp.NewResponse(fakehttp.WithErr(errors.New("connection refused"))))
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
}

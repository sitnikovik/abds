package invoices_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const twoInvoices = `{"invoices":[{"created_at":"2026-10-27 10:00:00"}],"count":2}`

func TestPagination_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the pages differ", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.
			NewPagination(
				fakehttp.NewSequence(
					fakehttp.NewResponse(
						fakehttp.WithStatus(http.StatusOK),
						fakehttp.WithBody(`{"invoices":[{"created_at":"2026-10-27 10:00:00"}],"count":2}`),
					),
					fakehttp.NewResponse(
						fakehttp.WithStatus(http.StatusOK),
						fakehttp.WithBody(`{"invoices":[{"created_at":"2026-10-27 10:00:00"}],"count":2}`),
					),
					fakehttp.NewResponse(
						fakehttp.WithStatus(http.StatusOK),
						fakehttp.WithBody(`{"invoices":[{"created_at":"2026-09-27 10:00:00"}],"count":2}`),
					),
				),
				"http://localhost:8080",
			).
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when the offset returns the same invoice", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(twoInvoices),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(twoInvoices),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(twoInvoices),
			),
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "offset must move"))
	})

	t.Run("fails when the second page is empty", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(twoInvoices),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(twoInvoices),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(`{"invoices":[],"count":2}`),
			),
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "second page"))
	})

	t.Run("fails when the default page is too large", func(t *testing.T) {
		t.Parallel()
		big := `{"invoices":[{"created_at":"1"},{"created_at":"2"},{"created_at":"3"},{"created_at":"4"},` +
			`{"created_at":"5"},{"created_at":"6"},{"created_at":"7"},{"created_at":"8"},{"created_at":"9"},` +
			`{"created_at":"10"},{"created_at":"11"}],"count":11}`
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(big),
			),
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "default page"))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusNotFound,
			`{"message":"x"}`,
			nil,
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			0,
			"",
			errors.New("connection refused"),
		)
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := invoices.NewPagination(cli, "http://localhost:8080").
			Run(
				context.Background(),
				vars,
			)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewPagination(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"http://localhost:8080",
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := invoices.NewPagination(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"http://localhost:8080",
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "invoices pagination"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

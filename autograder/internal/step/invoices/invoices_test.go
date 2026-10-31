package invoices_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithRealty returns variables that hold the realty id 1.
func VarsWithRealty() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	return vars
}

func TestList_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on a page with the total count", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				`{"invoices":[{"price":93000}],"count":36}`,
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when the page exceeds the limit", func(t *testing.T) {
		t.Parallel()
		body := `{"invoices":[{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1},{"price":1}],"count":36}`
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				body,
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.Contains(err, "page must not exceed the limit"))
	})

	t.Run("fails when the object has no invoices", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				`{"invoices":[],"count":0}`,
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.Contains(err, "object must have invoices"))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusNotFound,
				`{"message":"x"}`,
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails on a body that is not JSON", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				`{`,
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				0,
				"",
				errors.New("connection refused"),
			),
			baseURL,
		).Run(
			context.Background(),
			VarsWithRealty(),
		)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		err := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			baseURL,
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := invoices.NewList(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			baseURL,
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

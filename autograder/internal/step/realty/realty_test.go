package realty_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/realty"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithRealty returns variables that hold the realty id 1.
func VarsWithRealty() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	return vars
}

func TestGet_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the object with the expected fields", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":1,"address":"г. Москва","square":41,"type":"flat","floor":10,"porch":1}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("passes on a non-residential premises", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":1,"address":"г. Москва","square":50,"type":"nrcp","floor":10,"porch":1}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails on an unknown type", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":1,"address":"г. Москва","square":41,"type":"house"}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.Contains(err, "field type"))
	})

	t.Run("fails when the square is empty", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":1,"address":"г. Москва","square":0,"type":"flat","floor":10,"porch":1}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.Contains(err, "field square"))
	})

	t.Run("fails when the id differs", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":2,"address":"г. Москва","square":41,"type":"flat","floor":10,"porch":1}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.Contains(err, "field id"))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusNotFound,
				`{"message":"Not found"}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails on a body that is not JSON", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				`{`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				VarsWithRealty(),
			)
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
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
			)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		err := realty.NewGet(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := realty.NewGet(
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
				strings.Equal(s.Name(), "realty get"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestNotFound_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 404 with a message", func(t *testing.T) {
		t.Parallel()
		err := realty.NewNotFound(
			fakehttp.NewClient(
				http.StatusNotFound,
				`{"message":"Not found"}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when the object exists", func(t *testing.T) {
		t.Parallel()
		err := realty.NewNotFound(
			fakehttp.NewClient(
				http.StatusOK,
				`{"id":1}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "status code must be 404"))
	})

	t.Run("fails when the message is empty", func(t *testing.T) {
		t.Parallel()
		err := realty.NewNotFound(
			fakehttp.NewClient(
				http.StatusNotFound,
				`{"message":""}`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "field message"))
	})

	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		err := realty.NewNotFound(
			fakehttp.NewClient(
				http.StatusNotFound,
				`not json`,
				nil,
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		err := realty.NewNotFound(
			fakehttp.NewClient(
				0,
				"",
				errors.New("connection refused"),
			),
			baseURL,
		).
			Run(
				context.Background(),
				step.NewVariables(),
			)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := realty.NewNotFound(
			fakehttp.NewClient(
				http.StatusNotFound,
				"",
				nil,
			),
			baseURL,
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "realty get not found"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

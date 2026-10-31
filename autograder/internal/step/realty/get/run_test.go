package get_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/realty/get"
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

	t.Run("passes on the object with the expected fields", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"id":1,"address":"г. Москва","square":41,"type":"flat","floor":10,"porch":1}`,
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

	t.Run("passes on a non-residential premises", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"id":1,"address":"г. Москва","square":50,"type":"nrcp","floor":10,"porch":1}`,
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

	t.Run("fails on an unknown type", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"id":1,"address":"г. Москва","square":41,"type":"house"}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"field type",
			),
		)
	})

	t.Run("fails when the square is empty", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"id":1,"address":"г. Москва","square":0,"type":"flat","floor":10,"porch":1}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"field square",
			),
		)
	})

	t.Run("fails when the id differs", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							`{"id":2,"address":"г. Москва","square":41,"type":"flat","floor":10,"porch":1}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithRealty(),
					),
				"field id",
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
							`{"message":"Not found"}`,
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

	t.Run("fails on a body that is not JSON", func(t *testing.T) {
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
				strings.Equal(s.Name(), "realty get"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

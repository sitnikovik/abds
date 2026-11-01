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
	impl "github.com/sitnikovik/abds/autograder/internal/step/tariffs/get"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// validTariff is the body of the tariff in force.
const validTariff = `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,` +
	`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the tariff in force", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						fakehttp.NewClient(
							http.StatusOK,
							validTariff,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
			),
		)
	})

	t.Run("fails when the day price is empty", func(t *testing.T) {
		t.Parallel()
		var noerr error
		body := `{"id":1,"name":"День–Ночь","price_day_kwph":0,"price_night_kwph":310,` +
			`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
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
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				"field price_day_kwph",
			),
		)
	})

	t.Run("fails when the name is empty", func(t *testing.T) {
		t.Parallel()
		var noerr error
		body := `{"id":1,"name":"","price_day_kwph":620,"price_night_kwph":310,` +
			`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
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
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				"field name",
			),
		)
	})

	t.Run("fails when the validity is missing", func(t *testing.T) {
		t.Parallel()
		var noerr error
		body := `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,"valid_from":"","valid_to":""}`
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
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
					),
				"field valid_from",
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
							`{"message":"x"}`,
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						step.NewVariables(),
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
						step.NewVariables(),
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
						step.NewVariables(),
					),
				"connection refused",
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
				strings.Equal(s.Name(), "tariffs get"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

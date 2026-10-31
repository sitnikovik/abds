package endpoint_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/endpoint"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the expected status with a message", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set("user_id", "1")
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusNotFound,
				`{"message":"Not found"}`,
				nil,
			),
			baseURL,
			step.NewDesc("users get not found", "Checks 404."),
			endpoint.WithPath("/users/%s", "user_id"),
			endpoint.WithStatus(http.StatusNotFound),
		).Run(
			context.Background(),
			vars,
		)
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusOK,
				`{"message":"x"}`,
				nil,
			),
			baseURL,
			step.NewDesc("case", ""),
			endpoint.WithStatus(http.StatusBadRequest),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "status code must be 400"))
	})

	t.Run("fails when the message is empty", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusBadRequest,
				`{"message":""}`,
				nil,
			),
			baseURL,
			step.NewDesc("case", ""),
			endpoint.WithStatus(http.StatusBadRequest),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "field message"))
	})

	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusBadRequest,
				`x`,
				nil,
			),
			baseURL,
			step.NewDesc("case", ""),
			endpoint.WithStatus(http.StatusBadRequest),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				0,
				"",
				errors.New("connection refused"),
			),
			baseURL,
			step.NewDesc("case", ""),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when a path variable is not stored", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			baseURL,
			step.NewDesc("case", ""),
			endpoint.WithPath("/realty/%s", "realty_id"),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "realty_id"))
	})

	t.Run("fails on an invalid base url", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"://bad",
			step.NewDesc("case", ""),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, "build request"))
	})

	t.Run("uses the name and the description of the desc", func(t *testing.T) {
		t.Parallel()
		s := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			baseURL,
			step.NewDesc("tariffs bad request", "Checks 400."),
		)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "tariffs bad request"),
				strings.Equal(s.Description(), "Checks 400."),
			),
		)
	})

	t.Run("sends a body with the method option", func(t *testing.T) {
		t.Parallel()
		err := endpoint.NewStep(
			fakehttp.NewClient(
				http.StatusBadRequest,
				`{"message":"invalid \"gauge_id\" field"}`,
				nil,
			),
			baseURL,
			step.NewDesc("case", ""),
			endpoint.WithMethod(http.MethodPost),
			endpoint.WithBody(map[string]int{"gauge_id": 0}),
			endpoint.WithStatus(http.StatusBadRequest),
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.NoError(err))
	})
}

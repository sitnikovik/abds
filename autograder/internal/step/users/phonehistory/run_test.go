package phonehistory_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/users/phonehistory"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithUser returns variables that hold the user id 1.
func VarsWithUser() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.UserID, "1")
	return vars
}

// PhoneDB returns a fake database that stores the owner 1 with the phone +79000000001.
func PhoneDB() *fakeclickhouse.Client {
	return fakeclickhouse.NewClient(
		fakeclickhouse.WithStrings("dim_users", [][]string{{"+79000000001"}}),
	)
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the previous phone is kept in phone_prev", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(`{"phone":"+79999999999","phone_prev":"+79000000001"}`),
			),
		)
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
			),
		)
	})

	t.Run("fails when phone_prev does not keep the previous phone", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(`{"phone":"+79999999999","phone_prev":""}`),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"field phone_prev",
			),
		)
	})

	t.Run("fails when the phone is not changed", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(`{"phone":"+79000000001","phone_prev":"+79000000001"}`),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"field phone",
			),
		)
	})

	t.Run("fails on a wrong status code of the patch", func(t *testing.T) {
		t.Parallel()
		var noerr error
		cli := fakehttp.NewClient(
			http.StatusBadRequest,
			`{"message":"x"}`,
			noerr,
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"status code must be 204",
			),
		)
	})

	t.Run("fails when the read after the patch is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(`{`),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"decode body",
			),
		)
	})

	t.Run("fails on a wrong status code of the read", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusNotFound),
				fakehttp.WithBody(`{"message":"Not found"}`),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"status code must be 200",
			),
		)
	})

	t.Run("fails when the owner is not stored", func(t *testing.T) {
		t.Parallel()
		var rows [][]string
		var noerr error
		db := fakeclickhouse.NewClient(fakeclickhouse.WithStrings("dim_users", rows))
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						db,
						"energosbyt",
						fakehttp.NewClient(
							http.StatusOK,
							"",
							noerr,
						),
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"owner must be stored",
			),
		)
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			0,
			"",
			errors.New("connection refused"),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithUser(),
					),
				"connection refused",
			),
		)
	})

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						PhoneDB(),
						"energosbyt",
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
				fixture.UserID,
			),
		)
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		var noerr error
		s := impl.NewStep(
			PhoneDB(),
			"energosbyt",
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
				strings.Equal(s.Name(), "users phone history"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

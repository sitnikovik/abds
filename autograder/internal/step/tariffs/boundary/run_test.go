package boundary_test

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
	impl "github.com/sitnikovik/abds/autograder/internal/step/tariffs/boundary"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// oldVersion is the body of the first version of the tariff.
const oldVersion = `{"id":1,"name":"День–Ночь","price_day_kwph":450,"price_night_kwph":225,` +
	`"valid_from":"2023-01-01 00:00:00","valid_to":"2026-01-01 00:00:00"}`

// newVersion is the body of the second version of the tariff.
const newVersion = `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,` +
	`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`

// TwoVersions returns the stored starts of the two versions of the tariff.
func TwoVersions() *fakeclickhouse.Client {
	return fakeclickhouse.NewClient(
		fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}, {"2026-01-01 00:00:00"}}),
	)
}

// VarsWithTariff returns variables that hold the tariff id 1.
func VarsWithTariff() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.TariffID, "1")
	return vars
}

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the version changes on the boundary day", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(oldVersion),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(newVersion),
			),
		)
		assert.NoError(
			t,
			checkerr.NoError(
				impl.
					NewStep(
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
			),
		)
	})

	t.Run("fails when the version does not change on the boundary day", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(oldVersion),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(oldVersion),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"version must change",
			),
		)
	})

	t.Run("fails when the version on the day before does not cover that day", func(t *testing.T) {
		t.Parallel()
		later := `{"id":1,"name":"x","price_day_kwph":1,"price_night_kwph":1,` +
			`"valid_from":"2026-06-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(later),
			),
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(newVersion),
			),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"day before the boundary",
			),
		)
	})

	t.Run("fails when the tariff has one version only", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}}),
		)
		var noerr error
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
						VarsWithTariff(),
					),
				"two versions",
			),
		)
	})

	t.Run("fails when the start of the version is malformed", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}, {"not a date"}}),
		)
		var noerr error
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
						VarsWithTariff(),
					),
				"parse valid_from",
			),
		)
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
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
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"unexpected response",
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
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"connection refused",
			),
		)
	})

	t.Run("fails when the request for the boundary day fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(oldVersion),
			),
			fakehttp.NewResponse(fakehttp.WithErr(errors.New("connection reset"))),
		)
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"connection reset",
			),
		)
	})

	t.Run("fails when the version on the boundary day is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusOK),
				fakehttp.WithBody(oldVersion),
			),
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
						TwoVersions(),
						"energosbyt",
						cli,
						baseURL,
					).
					Run(
						context.Background(),
						VarsWithTariff(),
					),
				"decode body",
			),
		)
	})

	t.Run("fails when the tariff id is not stored", func(t *testing.T) {
		t.Parallel()
		var noerr error
		assert.NoError(
			t,
			checkerr.Contains(
				impl.
					NewStep(
						TwoVersions(),
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
				fixture.TariffID,
			),
		)
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		var noerr error
		s := impl.NewStep(
			TwoVersions(),
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
				strings.Equal(s.Name(), "tariffs boundary"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

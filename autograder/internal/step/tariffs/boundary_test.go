package tariffs_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/tariffs"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

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

const (
	oldVersion = `{"id":1,"name":"День–Ночь","price_day_kwph":450,"price_night_kwph":225,` +
		`"valid_from":"2023-01-01 00:00:00","valid_to":"2026-01-01 00:00:00"}`
	newVersion = `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,` +
		`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
)

func TestBoundary_Run(t *testing.T) {
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
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			cli,
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.NoError(err))
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
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			cli,
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.Contains(err, "version must change"))
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
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			cli,
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.Contains(err, "day before the boundary"))
	})

	t.Run("fails when the tariff has one version only", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}}),
		)
		err := tariffs.NewBoundary(
			db,
			"energosbyt",
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.Contains(err, "two versions"))
	})

	t.Run("fails when the start of the version is malformed", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}, {"not a date"}}),
		)
		err := tariffs.NewBoundary(
			db,
			"energosbyt",
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.Contains(err, "parse valid_from"))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(
				fakehttp.WithStatus(http.StatusNotFound),
				fakehttp.WithBody(`{"message":"Not found"}`),
			),
		)
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			cli,
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
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
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			cli,
			"http://localhost:8080",
		).Run(
			context.Background(),
			VarsWithTariff(),
		)
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the tariff id is not stored", func(t *testing.T) {
		t.Parallel()
		err := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
			fakehttp.NewClient(
				http.StatusOK,
				"",
				nil,
			),
			"http://localhost:8080",
		).Run(
			context.Background(),
			step.NewVariables(),
		)
		assert.NoError(t, checkerr.Contains(err, fixture.TariffID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBoundary(
			TwoVersions(),
			"energosbyt",
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
				strings.Equal(s.Name(), "tariffs boundary"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

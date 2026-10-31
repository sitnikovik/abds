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

// twoVersions returns the stored starts of the two versions of the tariff.
func twoVersions() *fakeclickhouse.Client {
	return fakeclickhouse.NewClient(
		fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}, {"2026-01-01 00:00:00"}}),
	)
}

// varsWithTariff returns variables that hold the tariff id 1.
func varsWithTariff() step.Variables {
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
			fakehttp.Canned{Status: http.StatusOK, Body: oldVersion},
			fakehttp.Canned{Status: http.StatusOK, Body: newVersion},
		)
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the version does not change on the boundary day", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusOK, Body: oldVersion},
			fakehttp.Canned{Status: http.StatusOK, Body: oldVersion},
		)
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "version must change"))
	})
	t.Run("fails when the version on the day before does not cover that day", func(t *testing.T) {
		t.Parallel()
		later := `{"id":1,"name":"x","price_day_kwph":1,"price_night_kwph":1,` +
			`"valid_from":"2026-06-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusOK, Body: later},
			fakehttp.Canned{Status: http.StatusOK, Body: newVersion},
		)
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "day before the boundary"))
	})
	t.Run("fails when the tariff has one version only", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}}),
		)
		s := tariffs.NewBoundary(db, "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "two versions"))
	})
	t.Run("fails when the start of the version is malformed", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(
			fakeclickhouse.WithStrings("dim_tariffs", [][]string{{"2023-01-01 00:00:00"}, {"not a date"}}),
		)
		s := tariffs.NewBoundary(db, "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "parse valid_from"))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNotFound, Body: `{"message":"Not found"}`},
		)
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(0, "", errors.New("connection refused"))
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the tariff id is not stored", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.TariffID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBoundary(twoVersions(), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "tariffs boundary"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

package readings_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/readings"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCorrection_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on one reading with the last sent values", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"day_kw":1250,"night_kw":300}],"count":1}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the date has two readings", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"day_kw":1250,"night_kw":300},{"day_kw":1,"night_kw":1}],"count":2}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "one actual reading"))
	})
	t.Run("fails when the day value is an old one", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"day_kw":1,"night_kw":300}],"count":1}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "field day_kw"))
	})
	t.Run("fails when the night value is an old one", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"day_kw":1250,"night_kw":1}],"count":1}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "field night_kw"))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(fakehttp.NewClient(http.StatusNotFound, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(fakehttp.NewClient(http.StatusOK, `{`, nil), baseURL)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := readings.NewCorrection(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings correction"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

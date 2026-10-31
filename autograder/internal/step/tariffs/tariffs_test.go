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
	"github.com/sitnikovik/abds/autograder/internal/step/tariffs"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

const validTariff = `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,` +
	`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`

func TestGet_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the tariff in force", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, validTariff, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the day price is empty", func(t *testing.T) {
		t.Parallel()
		body := `{"id":1,"name":"День–Ночь","price_day_kwph":0,"price_night_kwph":310,` +
			`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, body, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "field price_day_kwph"))
	})
	t.Run("fails when the name is empty", func(t *testing.T) {
		t.Parallel()
		body := `{"id":1,"name":"","price_day_kwph":620,"price_night_kwph":310,` +
			`"valid_from":"2026-01-01 00:00:00","valid_to":"2099-01-01 00:00:00"}`
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, body, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "field name"))
	})
	t.Run("fails when the validity is missing", func(t *testing.T) {
		t.Parallel()
		body := `{"id":1,"name":"День–Ночь","price_day_kwph":620,"price_night_kwph":310,"valid_from":"","valid_to":""}`
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, body, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "field valid_from"))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusNotFound, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails on a body that is not JSON", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, `{`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewGet(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewGet(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "tariffs get"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestBadRequest_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 400 with a message", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{"message":"invalid \"date\" parameter"}`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(http.StatusOK, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "status code must be 400"))
	})
	t.Run("fails when the message is empty", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{"message":""}`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "field message"))
	})
	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `x`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := tariffs.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "tariffs bad request"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

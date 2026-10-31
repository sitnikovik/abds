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

const baseURL = "http://localhost:8080"

// varsWithAll returns variables with every id the readings steps need.
func varsWithAll() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.UserID, "1")
	vars.Set(fixture.GaugeID, "1")
	vars.Set(fixture.TariffID, "1")
	return vars
}

func TestPost_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 201 with an empty body", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "status code must be 201"))
	})
	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "{}", nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "body must be empty"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("fails when a numeric id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings post"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestList_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the count matches the readings", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"reading_date":"2026-10-01"}],"count":1}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the count differs from the readings", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(
			fakehttp.NewClient(http.StatusOK, `{"readings":[{"reading_date":"2026-10-01"}],"count":2}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "count must match"))
	})
	t.Run("fails when October has no readings", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(
			fakehttp.NewClient(http.StatusOK, `{"readings":[],"count":0}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "October 2026"))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(fakehttp.NewClient(http.StatusNotFound, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails on a body that is not JSON", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(fakehttp.NewClient(http.StatusOK, `{`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := readings.NewList(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings list"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestBadRequest_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 400 with a message", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{"message":"invalid \"date_from\" parameter"}`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusOK, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "status code must be 400"))
	})
	t.Run("fails when the message is empty", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{"message":""}`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "field message"))
	})
	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `x`, nil), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), varsWithAll())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "readings list bad request"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

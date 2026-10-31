package gauges_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/gauges"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithRealtyAndGauge returns variables with a realty id and a new gauge id.
func VarsWithRealtyAndGauge() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.RealtyID, "1")
	vars.Set(fixture.NewGaugeID, "2")
	return vars
}

func TestPost_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 201 with an empty body", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), VarsWithRealtyAndGauge())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(http.StatusBadRequest, "", nil), baseURL)
		err := s.Run(context.Background(), VarsWithRealtyAndGauge())
		assert.NoError(t, checkerr.Contains(err, "status code must be 201"))
	})
	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(http.StatusCreated, "{}", nil), baseURL)
		err := s.Run(context.Background(), VarsWithRealtyAndGauge())
		assert.NoError(t, checkerr.Contains(err, "body must be empty"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), VarsWithRealtyAndGauge())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("fails when the gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		s := gauges.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.NewGaugeID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := gauges.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "gauges post"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

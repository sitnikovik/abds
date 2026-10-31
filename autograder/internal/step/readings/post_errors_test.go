package readings_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/readings"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestPost_RunErrors(t *testing.T) {
	t.Parallel()

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.GaugeID, "1")
		vars.Set(fixture.TariffID, "1")
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})

	t.Run("fails when the gauge id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.UserID, "1")
		vars.Set(fixture.TariffID, "1")
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.GaugeID))
	})

	t.Run("fails when the tariff id is not stored", func(t *testing.T) {
		t.Parallel()
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		vars.Set(fixture.UserID, "1")
		vars.Set(fixture.GaugeID, "1")
		s := readings.NewPost(fakehttp.NewClient(http.StatusCreated, "", nil), baseURL)
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, fixture.TariffID))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := readings.NewPost(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
}

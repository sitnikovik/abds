package readings_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/readings"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestBadRequest_RunErrors(t *testing.T) {
	t.Parallel()

	t.Run("fails when the realty id is not stored", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{"message":"x"}`, nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.RealtyID))
	})
	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		s := readings.NewBadRequest(fakehttp.NewClient(http.StatusBadRequest, `{`, nil), baseURL)
		err := s.Run(context.Background(), VarsWithAll())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
}

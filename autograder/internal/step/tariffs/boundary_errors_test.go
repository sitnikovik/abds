package tariffs_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step/tariffs"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestBoundary_RunErrors(t *testing.T) {
	t.Parallel()

	t.Run("fails when the version on the boundary day has a wrong status", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(oldVersion)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNotFound), fakehttp.WithBody(`{"message":"Not found"}`)),
		)
		s := tariffs.NewBoundary(TwoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the request for the boundary day fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(oldVersion)),
			fakehttp.NewResponse(fakehttp.WithErr(errors.New("connection reset"))),
		)
		s := tariffs.NewBoundary(TwoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "connection reset"))
	})
	t.Run("fails when the version on the boundary day is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(oldVersion)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(`{`)),
		)
		s := tariffs.NewBoundary(TwoVersions(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithTariff())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
}

package users_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step/users"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestFio_RunErrors(t *testing.T) {
	t.Parallel()

	t.Run("fails when the read after the patch has a wrong status", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Status: http.StatusNotFound, Body: `{"message":"Not found"}`},
		)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})
	t.Run("fails when the read after the patch is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Status: http.StatusOK, Body: `{`},
		)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
	t.Run("fails when the read after the patch fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Err: errors.New("connection reset")},
		)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "connection reset"))
	})
}

package users_test

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
	"github.com/sitnikovik/abds/autograder/internal/step/users"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

// PhoneDB returns a fake database that stores the owner 1 with the phone +79000000001.
func PhoneDB() *fakeclickhouse.Client {
	return fakeclickhouse.NewClient(
		fakeclickhouse.WithStrings("dim_users", [][]string{{"+79000000001"}}),
	)
}

func TestPhoneHistory_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the previous phone is kept in phone_prev", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(`{"phone":"+79999999999","phone_prev":"+79000000001"}`)),
		)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when phone_prev does not keep the previous phone", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(`{"phone":"+79999999999","phone_prev":""}`)),
		)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "field phone_prev"))
	})
	t.Run("fails when the phone is not changed", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(`{"phone":"+79000000001","phone_prev":"+79000000001"}`)),
		)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "field phone"))
	})
	t.Run("fails on a wrong status code of the patch", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(http.StatusBadRequest, `{"message":"x"}`, nil)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "status code must be 204"))
	})
	t.Run("fails when the read after the patch is not JSON", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusOK), fakehttp.WithBody(`{`)),
		)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})
	t.Run("fails on a wrong status code of the read", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNoContent)),
			fakehttp.NewResponse(fakehttp.WithStatus(http.StatusNotFound), fakehttp.WithBody(`{"message":"Not found"}`)),
		)
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "status code must be 200"))
	})
	t.Run("fails when the owner is not stored", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(fakeclickhouse.WithStrings("dim_users", nil))
		s := users.NewPhoneHistory(db, "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "owner must be stored"))
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(0, "", errors.New("connection refused"))
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := users.NewPhoneHistory(PhoneDB(), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users phone history"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

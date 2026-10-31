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

// ownerDB returns a fake database that stores the owner 1 with the full name Иванов Иван and no other names.
func ownerDB(stored int) *fakeclickhouse.Client {
	return fakeclickhouse.NewClient(
		fakeclickhouse.WithStrings("dim_users", [][]string{{"Иванов Иван"}}),
		fakeclickhouse.WithCount("fio = ?", uint64(stored)),
	)
}

func TestFio_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes when the new name is returned and the previous one is gone", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Status: http.StatusOK, Body: `{"fio":"Автограйдер Тест"}`},
		)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when the returned name differs", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Status: http.StatusOK, Body: `{"fio":"Иванов Иван"}`},
		)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "field fio"))
	})
	t.Run("fails when the previous name is still stored", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewSequence(
			fakehttp.Canned{Status: http.StatusNoContent},
			fakehttp.Canned{Status: http.StatusOK, Body: `{"fio":"Автограйдер Тест"}`},
		)
		s := users.NewFio(ownerDB(1), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "previous name must not be stored"))
	})
	t.Run("fails when the owner is not stored", func(t *testing.T) {
		t.Parallel()
		db := fakeclickhouse.NewClient(fakeclickhouse.WithStrings("dim_users", nil))
		s := users.NewFio(db, "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "owner must be stored"))
	})
	t.Run("fails on a wrong status code of the patch", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(http.StatusBadRequest, `{"message":"x"}`, nil)
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "status code must be 204"))
	})
	t.Run("fails when the patch request fails", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(0, "", errors.New("connection refused"))
		s := users.NewFio(ownerDB(0), "energosbyt", cli, "http://localhost:8080")
		err := s.Run(context.Background(), varsWithUser())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})
	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		s := users.NewFio(ownerDB(0), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := users.NewFio(ownerDB(0), "energosbyt", fakehttp.NewClient(http.StatusOK, "", nil), "http://localhost:8080")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users fio"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

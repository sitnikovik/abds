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
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

const baseURL = "http://localhost:8080"

// VarsWithUser returns variables that hold the user id 1.
func VarsWithUser() step.Variables {
	vars := step.NewVariables()
	vars.Set(fixture.UserID, "1")
	return vars
}

func TestGet_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the owner with the expected fields", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(
			fakehttp.NewClient(http.StatusOK, `{"id":1,"fio":"Иванов Иван","phone":"+79000000001","phone_prev":""}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when the id differs", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(
			fakehttp.NewClient(http.StatusOK, `{"id":2,"fio":"Иванов Иван","phone":"+79000000001"}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "field id"))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(
			fakehttp.NewClient(http.StatusNotFound, `{"message":"Not found"}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails on a body that is not JSON", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(
			fakehttp.NewClient(http.StatusOK, `{`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "decode body"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(
			fakehttp.NewClient(0, "", errors.New("connection refused")),
			baseURL,
		)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := users.NewGet(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users get"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestNotFound_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 404 with a message", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(
			fakehttp.NewClient(http.StatusNotFound, `{"message":"Not found"}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails when the owner exists", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(
			fakehttp.NewClient(http.StatusOK, `{"id":1}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "status code must be 404"))
	})

	t.Run("fails when the message is empty", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(
			fakehttp.NewClient(http.StatusNotFound, `{"message":""}`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "field message"))
	})

	t.Run("fails when the body is not JSON", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(
			fakehttp.NewClient(http.StatusNotFound, `not json`, nil),
			baseURL,
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "unexpected response"))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(
			fakehttp.NewClient(0, "", errors.New("connection refused")),
			baseURL,
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := users.NewNotFound(fakehttp.NewClient(http.StatusNotFound, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users get not found"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

func TestPatch_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 204 with an empty body", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(http.StatusNoContent, "", nil), baseURL)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.NoError(err))
	})

	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(http.StatusOK, "", nil), baseURL)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "status code must be 204"))
	})

	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(http.StatusNoContent, "x", nil), baseURL)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Error(err))
	})

	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(0, "", errors.New("connection refused")), baseURL)
		err := s.Run(context.Background(), VarsWithUser())
		assert.NoError(t, checkerr.Contains(err, "connection refused"))
	})

	t.Run("fails when the user id is not stored", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(http.StatusNoContent, "", nil), baseURL)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(t, checkerr.Contains(err, fixture.UserID))
	})

	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := users.NewPatch(fakehttp.NewClient(http.StatusNoContent, "", nil), baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "users patch"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

package ping_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	errs "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/ping"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on 200 with an empty body", func(t *testing.T) {
		t.Parallel()
		s := ping.NewStep(
			httpclient.NewClient(http.StatusOK, "", nil),
			"http://localhost:8080/ping",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			errs.NoError(err),
		)
	})
	t.Run("fails on a wrong status code", func(t *testing.T) {
		t.Parallel()
		s := ping.NewStep(
			httpclient.NewClient(http.StatusInternalServerError, "", nil),
			"http://localhost:8080/ping",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			errs.Error(err),
		)
	})
	t.Run("fails on a non-empty body", func(t *testing.T) {
		t.Parallel()
		s := ping.NewStep(
			httpclient.NewClient(http.StatusOK, `{"status":"ok"}`, nil),
			"http://localhost:8080/ping",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			errs.Error(err),
		)
	})
	t.Run("fails when the request fails", func(t *testing.T) {
		t.Parallel()
		s := ping.NewStep(
			httpclient.NewClient(0, "", errors.New("connection refused")),
			"http://localhost:8080/ping",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			expect.All(
				errs.Error(err),
				errs.Contains(err, "connection refused"),
			),
		)
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := ping.NewStep(httpclient.NewClient(http.StatusOK, "", nil), "")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "ping"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

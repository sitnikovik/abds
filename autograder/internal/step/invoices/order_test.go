package invoices_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	checkerr "github.com/sitnikovik/buggles/expect/errors"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestList_Order(t *testing.T) {
	t.Parallel()

	t.Run("passes when invoices are sorted newest first", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusOK,
			`{"invoices":[{"created_at":"2026-10-27 10:00:00"},{"created_at":"2026-09-27 10:00:00"}],"count":2}`,
			nil,
		)
		s := invoices.NewList(cli, "http://localhost:8080")
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.NoError(err))
	})
	t.Run("fails when invoices are sorted oldest first", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusOK,
			`{"invoices":[{"created_at":"2026-09-27 10:00:00"},{"created_at":"2026-10-27 10:00:00"}],"count":2}`,
			nil,
		)
		s := invoices.NewList(cli, "http://localhost:8080")
		vars := step.NewVariables()
		vars.Set(fixture.RealtyID, "1")
		err := s.Run(context.Background(), vars)
		assert.NoError(t, checkerr.Contains(err, "sorted by created_at"))
	})
}

package invoices_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step/invoices"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCases_Names(t *testing.T) {
	t.Parallel()

	t.Run("names the error cases", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusBadRequest,
			"",
			nil,
		)
		bad := invoices.NewBadRequest(cli, "http://localhost:8080")
		notFound := invoices.NewNotFound(cli, "http://localhost:8080")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(bad.Name(), "invoices bad request"),
				strings.Equal(notFound.Name(), "invoices not found"),
				strings.NotEmpty(notFound.Description()),
			),
		)
	})
}

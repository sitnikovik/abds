package tariffs_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step/tariffs"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCases_Names(t *testing.T) {
	t.Parallel()

	t.Run("names the error cases", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(http.StatusBadRequest, "", nil)
		invalid := tariffs.NewInvalidDate(cli, baseURL)
		notFound := tariffs.NewNotFound(cli, baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(invalid.Name(), "tariffs invalid date"),
				strings.Equal(notFound.Name(), "tariffs not found"),
				strings.NotEmpty(notFound.Description()),
			),
		)
	})
}

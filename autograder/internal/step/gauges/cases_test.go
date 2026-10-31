package gauges_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step/gauges"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCases_Names(t *testing.T) {
	t.Parallel()

	t.Run("names the post error cases", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(http.StatusBadRequest, "", nil)
		bad := gauges.NewPostBadRequest(cli, baseURL)
		notFound := gauges.NewPostNotFound(cli, baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(bad.Name(), "gauges post bad request"),
				strings.Equal(notFound.Name(), "gauges post not found"),
				strings.NotEmpty(notFound.Description()),
			),
		)
	})
}

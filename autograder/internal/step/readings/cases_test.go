package readings_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step/readings"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCases_Names(t *testing.T) {
	t.Parallel()

	t.Run("names the list and post error cases", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			http.StatusNotFound,
			"",
			nil,
		)
		listNotFound := readings.NewListNotFound(cli, baseURL)
		postBad := readings.NewPostBadRequest(cli, baseURL)
		postNotFound := readings.NewPostNotFound(cli, baseURL)
		assert.NoError(
			t,
			expect.All(
				strings.Equal(listNotFound.Name(), "readings list not found"),
				strings.Equal(postBad.Name(), "readings post bad request"),
				strings.Equal(postNotFound.Name(), "readings post not found"),
				strings.NotEmpty(postNotFound.Description()),
			),
		)
	})
}

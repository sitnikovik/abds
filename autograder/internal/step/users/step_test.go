package users_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step/users"
	fakehttp "github.com/sitnikovik/abds/autograder/internal/test/fake/httpclient"
)

func TestCases_Names(t *testing.T) {
	t.Parallel()

	t.Run("names the patch error cases", func(t *testing.T) {
		t.Parallel()
		cli := fakehttp.NewClient(
			400,
			"",
			nil,
		)
		bad := users.NewPatchBadRequest(cli, "http://localhost:8080")
		notFound := users.NewPatchNotFound(cli, "http://localhost:8080")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(bad.Name(), "users patch bad request"),
				strings.NotEmpty(bad.Description()),
				strings.Equal(notFound.Name(), "users patch not found"),
				strings.NotEmpty(notFound.Description()),
			),
		)
	})
}

package lab1_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/autograder/lab1"
	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/driver"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/env"
)

func TestSteps(t *testing.T) {
	t.Parallel()

	t.Run("returns the checks in the order they run", func(t *testing.T) {
		t.Parallel()
		cfg, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "8080",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(t, checkerr.NoError(err))
		steps := lab1.Steps(
			cfg,
			clickhouse.NewClient(driver.NewConn(nil, nil)),
			&http.Client{},
		)
		assert.NoError(
			t,
			expect.All(
				number.Equal(len(steps), 33),
				strings.Equal(steps[0].Name(), "ping"),
				strings.Equal(steps[1].Name(), "schema"),
				strings.Equal(steps[2].Name(), "data"),
				strings.Equal(steps[3].Name(), "fixture"),
				strings.Equal(steps[len(steps)-1].Name(), "gauges replacement"),
			),
		)
	})
}

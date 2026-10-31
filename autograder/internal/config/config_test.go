package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/test/fake/env"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("reads all settings", func(t *testing.T) {
		t.Parallel()
		cfg, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "8080",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
			"CLICKHOUSE_PASSWORD": "secret",
		}))
		assert.NoError(
			t,
			expect.All(
				errors.NoError(err),
				strings.Equal(cfg.AppPort(), "8080"),
				strings.Equal(cfg.ClickHouseAddr(), "localhost:9000"),
				strings.Equal(cfg.ClickHouseDatabase(), "energosbyt"),
				strings.Equal(cfg.ClickHouseUser(), "default"),
				strings.Equal(cfg.ClickHousePassword(), "secret"),
			),
		)
	})
	t.Run("allows empty password", func(t *testing.T) {
		t.Parallel()
		cfg, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "8080",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(
			t,
			expect.All(
				errors.NoError(err),
				strings.Empty(cfg.ClickHousePassword()),
			),
		)
	})
	t.Run("fails when a required variable is missing", func(t *testing.T) {
		t.Parallel()
		_, err := config.Load(env.NewLookup(map[string]string{
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(
			t,
			expect.All(
				errors.Error(err),
				errors.Contains(err, "APP_PORT"),
			),
		)
	})
	t.Run("fails when a required variable is empty", func(t *testing.T) {
		t.Parallel()
		_, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(
			t,
			errors.Error(err),
		)
	})
	t.Run("fails when a port is not a number", func(t *testing.T) {
		t.Parallel()
		_, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "eighty",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "9000",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(
			t,
			expect.All(
				errors.Error(err),
				errors.Contains(err, "APP_PORT"),
			),
		)
	})
	t.Run("fails when a port is zero", func(t *testing.T) {
		t.Parallel()
		_, err := config.Load(env.NewLookup(map[string]string{
			"APP_PORT":            "8080",
			"CLICKHOUSE_HOST":     "localhost",
			"CLICKHOUSE_PORT":     "0",
			"CLICKHOUSE_DATABASE": "energosbyt",
			"CLICKHOUSE_USER":     "default",
		}))
		assert.NoError(
			t,
			errors.Contains(err, "CLICKHOUSE_PORT"),
		)
	})
}

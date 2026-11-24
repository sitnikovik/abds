package clickhouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/clickhouse"
	"energosbyt/internal/config/clickhouse/auth"
	"energosbyt/internal/config/clickhouse/connection"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with all fields", func(t *testing.T) {
		t.Parallel()
		conn := connection.NewConfig(
			"localhost",
			9000,
			"test",
		)
		authCfg := auth.NewConfig(
			"admin",
			"secret",
		)
		cfg := clickhouse.NewConfig(
			conn,
			authCfg,
		)
		assert.Equal(t, "localhost", cfg.Connection().Host())
	})

	t.Run("preserves auth credentials", func(t *testing.T) {
		t.Parallel()
		conn := connection.NewConfig(
			"localhost",
			9000,
			"test",
		)
		authCfg := auth.NewConfig(
			"user123",
			"pass456",
		)
		cfg := clickhouse.NewConfig(
			conn,
			authCfg,
		)
		assert.Equal(t, "user123", cfg.Auth().Username())
	})

	t.Run("returns connection config", func(t *testing.T) {
		t.Parallel()
		conn := connection.NewConfig(
			"ch-server",
			8123,
			"production",
		)
		cfg := clickhouse.NewConfig(
			conn,
			auth.NewConfig(
				"",
				"",
			),
		)
		assert.Equal(t, 8123, cfg.Connection().Port())
	})

	t.Run("returns auth config", func(t *testing.T) {
		t.Parallel()
		authCfg := auth.NewConfig(
			"root",
			"rootpass",
		)
		cfg := clickhouse.NewConfig(
			connection.NewConfig(
				"",
				0,
				"",
			),
			authCfg,
		)
		assert.Equal(t, "rootpass", cfg.Auth().Password())
	})
}

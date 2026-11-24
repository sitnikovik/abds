package connection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/clickhouse/connection"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with all fields", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig("localhost", 9000, "energosbyt")
		assert.Equal(t, "localhost", cfg.Host())
	})

	t.Run("preserves port value", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig("ch-server", 8123, "test")
		assert.Equal(t, 8123, cfg.Port())
	})

	t.Run("preserves database name", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig("localhost", 9000, "production")
		assert.Equal(t, "production", cfg.Database())
	})

	t.Run("handles empty values", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig("", 0, "")
		assert.Equal(t, "", cfg.Host())
	})
}

package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/clickhouse/auth"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with credentials", func(t *testing.T) {
		t.Parallel()
		cfg := auth.NewConfig("admin", "secret")
		assert.Equal(t, "admin", cfg.Username())
	})

	t.Run("preserves password", func(t *testing.T) {
		t.Parallel()
		cfg := auth.NewConfig("user", "pass123")
		assert.Equal(t, "pass123", cfg.Password())
	})

	t.Run("handles empty credentials", func(t *testing.T) {
		t.Parallel()
		cfg := auth.NewConfig("", "")
		assert.Equal(t, "", cfg.Username())
	})
}

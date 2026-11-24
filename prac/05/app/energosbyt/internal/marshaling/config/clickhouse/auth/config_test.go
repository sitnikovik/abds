package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	marshal "energosbyt/internal/marshaling/config/clickhouse/auth"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Username: "admin",
			Password: "secret",
		}
		domain := cfg.Config()
		assert.Equal(t, "admin", domain.Username())
	})

	t.Run("preserves password", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Username: "user",
			Password: "pass123",
		}
		domain := cfg.Config()
		assert.Equal(t, "pass123", domain.Password())
	})

	t.Run("handles empty credentials", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Username: "",
			Password: "",
		}
		domain := cfg.Config()
		assert.Equal(t, "", domain.Username())
	})
}

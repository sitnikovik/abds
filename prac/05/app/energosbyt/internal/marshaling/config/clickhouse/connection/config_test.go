package connection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	marshal "energosbyt/internal/marshaling/config/clickhouse/connection"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Host:     "localhost",
			Port:     9000,
			Database: "energosbyt",
		}
		domain := cfg.Config()
		assert.Equal(t, "localhost", domain.Host())
	})

	t.Run("preserves port value", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Host:     "ch-server",
			Port:     8123,
			Database: "test",
		}
		domain := cfg.Config()
		assert.Equal(t, 8123, domain.Port())
	})

	t.Run("preserves database name", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Host:     "localhost",
			Port:     9000,
			Database: "production",
		}
		domain := cfg.Config()
		assert.Equal(t, "production", domain.Database())
	})

	t.Run("handles empty values", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Host:     "",
			Port:     0,
			Database: "",
		}
		domain := cfg.Config()
		assert.Equal(t, "", domain.Host())
	})
}

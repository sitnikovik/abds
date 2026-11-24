package clickhouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	marshal "energosbyt/internal/marshaling/config/clickhouse"
	auth "energosbyt/internal/marshaling/config/clickhouse/auth"
	connection "energosbyt/internal/marshaling/config/clickhouse/connection"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Connection: connection.Config{
				Host:     "localhost",
				Port:     9000,
				Database: "test",
			},
			Auth: auth.Config{
				Username: "admin",
				Password: "secret",
			},
		}
		domain := cfg.Config()
		assert.Equal(t, "localhost", domain.Connection().Host())
	})

	t.Run("preserves auth credentials", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Connection: connection.Config{
				Host:     "localhost",
				Port:     9000,
				Database: "test",
			},
			Auth: auth.Config{
				Username: "user123",
				Password: "pass456",
			},
		}
		domain := cfg.Config()
		assert.Equal(t, "user123", domain.Auth().Username())
	})
}

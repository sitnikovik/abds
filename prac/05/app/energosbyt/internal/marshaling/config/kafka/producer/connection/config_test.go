package connection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	marshal "energosbyt/internal/marshaling/config/kafka/producer/connection"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Timeout: 5000,
		}
		domain := cfg.Config()
		assert.Equal(t, 5000, domain.Timeout())
	})

	t.Run("preserves timeout value", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Timeout: 10000,
		}
		domain := cfg.Config()
		assert.Equal(t, 10000, domain.Timeout())
	})

	t.Run("handles zero timeout", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Timeout: 0,
		}
		domain := cfg.Config()
		assert.Equal(t, 0, domain.Timeout())
	})
}

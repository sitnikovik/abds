package connection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka/producer/connection"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with timeout", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig(5000)
		assert.Equal(t, 5000, cfg.Timeout())
	})

	t.Run("creates config with zero timeout", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig(0)
		assert.Equal(t, 0, cfg.Timeout())
	})

	t.Run("preserves large timeout value", func(t *testing.T) {
		t.Parallel()
		cfg := connection.NewConfig(30000)
		assert.Equal(t, 30000, cfg.Timeout())
	})
}

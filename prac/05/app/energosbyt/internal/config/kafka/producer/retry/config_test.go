package retry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka/producer/retry"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with max retries", func(t *testing.T) {
		t.Parallel()
		cfg := retry.NewConfig(3)
		assert.Equal(t, 3, cfg.Max())
	})

	t.Run("creates config with zero retries", func(t *testing.T) {
		t.Parallel()
		cfg := retry.NewConfig(0)
		assert.Equal(t, 0, cfg.Max())
	})

	t.Run("preserves large retry value", func(t *testing.T) {
		t.Parallel()
		cfg := retry.NewConfig(10)
		assert.Equal(t, 10, cfg.Max())
	})
}

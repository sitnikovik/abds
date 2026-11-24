package retry_test

import (
	"testing"

	marshal "energosbyt/internal/marshaling/config/kafka/producer/retry"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Max: 3,
		}
		_ = cfg.Config()
	})
}

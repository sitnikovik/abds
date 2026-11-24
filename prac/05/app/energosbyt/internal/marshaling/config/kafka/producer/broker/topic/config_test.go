package topic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	marshal "energosbyt/internal/marshaling/config/kafka/producer/broker/topic"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Name: "meter_readings",
		}
		domain := cfg.Config()
		assert.Equal(t, "meter_readings", domain.Name())
	})

	t.Run("preserves topic name", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Name: "energy_events",
		}
		domain := cfg.Config()
		assert.Equal(t, "energy_events", domain.Name())
	})

	t.Run("handles empty name", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Name: "",
		}
		domain := cfg.Config()
		assert.Equal(t, "", domain.Name())
	})
}

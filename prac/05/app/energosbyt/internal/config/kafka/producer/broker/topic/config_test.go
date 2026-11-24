package topic_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka/producer/broker/topic"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with name", func(t *testing.T) {
		t.Parallel()
		cfg := topic.NewConfig("meter_readings")
		assert.Equal(t, "meter_readings", cfg.Name())
	})

	t.Run("creates config with empty name", func(t *testing.T) {
		t.Parallel()
		cfg := topic.NewConfig("")
		assert.Equal(t, "", cfg.Name())
	})

	t.Run("preserves complex topic name", func(t *testing.T) {
		t.Parallel()
		cfg := topic.NewConfig("energosbyt.events.2024")
		assert.Equal(t, "energosbyt.events.2024", cfg.Name())
	})
}

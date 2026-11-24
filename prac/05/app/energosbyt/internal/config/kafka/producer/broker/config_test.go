package broker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka/producer/broker"
	"energosbyt/internal/config/kafka/producer/broker/topic"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with all fields", func(t *testing.T) {
		t.Parallel()
		topics := []topic.Config{
			topic.NewConfig("test"),
		}
		cfg := broker.NewConfig("localhost", 9092, topics)
		assert.Equal(t, "localhost", cfg.Name())
	})

	t.Run("preserves port value", func(t *testing.T) {
		t.Parallel()
		cfg := broker.NewConfig("broker1", 9093, nil)
		assert.Equal(t, 9093, cfg.Port())
	})

	t.Run("preserves topics list", func(t *testing.T) {
		t.Parallel()
		topics := []topic.Config{
			topic.NewConfig("topic1"),
			topic.NewConfig("topic2"),
		}
		cfg := broker.NewConfig("broker1", 9092, topics)
		assert.Equal(t, 2, len(cfg.Topics()))
	})

	t.Run("handles empty topics", func(t *testing.T) {
		t.Parallel()
		cfg := broker.NewConfig("broker1", 9092, nil)
		assert.Equal(t, 0, len(cfg.Topics()))
	})

	t.Run("returns string representation", func(t *testing.T) {
		t.Parallel()
		cfg := broker.NewConfig("localhost", 9092, nil)
		assert.Equal(t, "localhost:9092", cfg.String())
	})
}

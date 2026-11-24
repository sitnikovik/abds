package broker_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	broker "energosbyt/internal/marshaling/config/kafka/producer/broker"
	topic "energosbyt/internal/marshaling/config/kafka/producer/broker/topic"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := broker.Config{
			Name: "localhost",
			Port: 9092,
			Topics: []topic.Config{
				{Name: "test"},
			},
		}
		domain := cfg.Config()
		assert.Equal(t, "localhost", domain.Name())
	})

	t.Run("preserves port value", func(t *testing.T) {
		t.Parallel()
		cfg := broker.Config{
			Name: "broker1",
			Port: 9093,
		}
		domain := cfg.Config()
		assert.Equal(t, 9093, domain.Port())
	})

	t.Run("preserves topics count", func(t *testing.T) {
		t.Parallel()
		cfg := broker.Config{
			Name: "broker1",
			Port: 9092,
			Topics: []topic.Config{
				{Name: "topic1"},
				{Name: "topic2"},
			},
		}
		domain := cfg.Config()
		assert.Equal(t, 2, len(domain.Topics()))
	})

	t.Run("handles empty topics", func(t *testing.T) {
		t.Parallel()
		cfg := broker.Config{
			Name: "broker1",
			Port: 9092,
		}
		domain := cfg.Config()
		assert.Equal(t, 0, len(domain.Topics()))
	})
}

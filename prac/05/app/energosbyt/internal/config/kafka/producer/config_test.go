package producer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka/producer"
	"energosbyt/internal/config/kafka/producer/broker"
	"energosbyt/internal/config/kafka/producer/connection"
	"energosbyt/internal/config/kafka/producer/retry"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with all fields", func(t *testing.T) {
		t.Parallel()
		brokers := []broker.Config{
			broker.NewConfig(
				"localhost",
				9092,
				nil,
			),
		}
		conn := connection.NewConfig(5000)
		retryCfg := retry.NewConfig(3)
		cfg := producer.NewConfig(
			brokers,
			conn,
			retryCfg,
			1,
		)
		assert.Equal(t, 1, len(cfg.Brokers()))
	})

	t.Run("preserves ack value", func(t *testing.T) {
		t.Parallel()
		cfg := producer.NewConfig(
			nil,
			connection.NewConfig(5000),
			retry.NewConfig(3),
			-1,
		)
		assert.Equal(t, -1, cfg.Ack())
	})

	t.Run("preserves brokers list", func(t *testing.T) {
		t.Parallel()
		brokers := []broker.Config{
			broker.NewConfig(
				"broker1",
				9092,
				nil,
			),
			broker.NewConfig(
				"broker2",
				9093,
				nil,
			),
		}
		cfg := producer.NewConfig(
			brokers,
			connection.NewConfig(5000),
			retry.NewConfig(3),
			1,
		)
		assert.Equal(t, 2, len(cfg.Brokers()))
	})

	t.Run("handles empty brokers", func(t *testing.T) {
		t.Parallel()
		cfg := producer.NewConfig(
			nil,
			connection.NewConfig(5000),
			retry.NewConfig(3),
			1,
		)
		assert.Equal(t, 0, len(cfg.Brokers()))
	})

	t.Run("returns connection config", func(t *testing.T) {
		t.Parallel()
		conn := connection.NewConfig(10000)
		cfg := producer.NewConfig(
			nil,
			conn,
			retry.NewConfig(3),
			1,
		)
		assert.Equal(t, 10000, cfg.Connection().Timeout())
	})

	t.Run("returns retry config", func(t *testing.T) {
		t.Parallel()
		retryCfg := retry.NewConfig(5)
		cfg := producer.NewConfig(
			nil,
			connection.NewConfig(5000),
			retryCfg,
			1,
		)
		assert.Equal(t, 5, cfg.Retry().Max())
	})
}

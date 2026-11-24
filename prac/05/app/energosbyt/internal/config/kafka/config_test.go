package kafka_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/kafka"
	"energosbyt/internal/config/kafka/producer"
	"energosbyt/internal/config/kafka/producer/broker"
	"energosbyt/internal/config/kafka/producer/connection"
	"energosbyt/internal/config/kafka/producer/retry"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with producer", func(t *testing.T) {
		t.Parallel()
		prodCfg := producer.NewConfig(
			[]broker.Config{
				broker.NewConfig(
					"localhost",
					9092,
					nil,
				),
			},
			connection.NewConfig(5000),
			retry.NewConfig(3),
			1,
		)
		cfg := kafka.NewConfig(prodCfg)
		assert.Equal(t, 1, len(cfg.Producer().Brokers()))
	})

	t.Run("returns producer config", func(t *testing.T) {
		t.Parallel()
		prodCfg := producer.NewConfig(
			nil,
			connection.NewConfig(10000),
			retry.NewConfig(5),
			-1,
		)
		cfg := kafka.NewConfig(prodCfg)
		assert.Equal(t, -1, cfg.Producer().Ack())
	})
}

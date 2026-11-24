package producer_test

import (
	"testing"

	marshal "energosbyt/internal/marshaling/config/kafka/producer"
	broker "energosbyt/internal/marshaling/config/kafka/producer/broker"
	connection "energosbyt/internal/marshaling/config/kafka/producer/connection"
	retry "energosbyt/internal/marshaling/config/kafka/producer/retry"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Brokers: []broker.Config{
				{Name: "localhost", Port: 9092},
			},
			Connection: connection.Config{
				Timeout: 5000,
			},
			Retry: retry.Config{
				Max: 3,
			},
			Ack: 1,
		}
		_ = cfg.Config()
	})
}

package config

import "abds-producer/internal/infra/client/kafka/config/producer"

// Kafka описывает параметры конфигурации Kafka.
type Kafka struct {
	// Producer - параметры продюсера Kafka.
	Producer producer.Producer
}

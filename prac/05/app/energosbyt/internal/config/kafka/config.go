package kafka

import "energosbyt/internal/config/kafka/producer"

// Config holds the complete Kafka configuration including producer settings.
type Config struct {
	// producer stores the producer configuration for sending messages.
	producer producer.Config
}

// NewConfig creates a new Config with the given producer configuration.
func NewConfig(producer producer.Config) Config {
	return Config{
		producer: producer,
	}
}

// Producer returns the producer configuration.
func (c Config) Producer() producer.Config {
	return c.producer
}

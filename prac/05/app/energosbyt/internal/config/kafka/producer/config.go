package producer

import (
	"energosbyt/internal/config/kafka/producer/broker"
	"energosbyt/internal/config/kafka/producer/connection"
	"energosbyt/internal/config/kafka/producer/retry"
)

// Config holds the complete Kafka producer configuration.
type Config struct {
	// brokers stores the list of broker configurations.
	brokers []broker.Config

	// connection stores the connection settings for the producer.
	connection connection.Config

	// retry stores the retry configuration for failed operations.
	retry retry.Config

	// ack stores the required number of acknowledgements from brokers.
	ack int
}

// NewConfig creates a new Config with the given producer parameters.
func NewConfig(
	brokers []broker.Config,
	connection connection.Config,
	retry retry.Config,
	ack int,
) Config {
	return Config{
		brokers:    brokers,
		connection: connection,
		retry:      retry,
		ack:        ack,
	}
}

// Brokers returns the list of broker configurations.
func (c Config) Brokers() []broker.Config {
	return c.brokers
}

// Connection returns the connection configuration.
func (c Config) Connection() connection.Config {
	return c.connection
}

// Retry returns the retry configuration.
func (c Config) Retry() retry.Config {
	return c.retry
}

// Ack returns the required number of acknowledgements from brokers.
func (c Config) Ack() int {
	return c.ack
}

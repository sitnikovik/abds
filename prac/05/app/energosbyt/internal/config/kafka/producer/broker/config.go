package broker

import (
	"fmt"

	"energosbyt/internal/config/kafka/producer/broker/topic"
)

// Config holds Kafka broker configuration including connection details and topics.
type Config struct {
	// name stores the hostname or IP address of the broker.
	name string

	// port stores the port number of the broker.
	port int

	// topics stores the list of topic configurations for this broker.
	topics []topic.Config
}

// NewConfig creates a new Config with the given broker parameters.
func NewConfig(
	name string,
	port int,
	topics []topic.Config,
) Config {
	return Config{
		name:   name,
		port:   port,
		topics: topics,
	}
}

// Name returns the broker hostname or IP address.
func (c Config) Name() string {
	return c.name
}

// Port returns the broker port number.
func (c Config) Port() int {
	return c.port
}

// Topics returns the list of topic configurations for this broker.
func (c Config) Topics() []topic.Config {
	return c.topics
}

// String returns the broker address as a string in host:port format.
func (c Config) String() string {
	return fmt.Sprintf("%s:%d", c.Name(), c.Port())
}

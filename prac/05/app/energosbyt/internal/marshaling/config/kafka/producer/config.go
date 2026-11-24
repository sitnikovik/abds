package producer

import (
	"energosbyt/internal/config/kafka/producer"
	domain "energosbyt/internal/config/kafka/producer/broker"
	broker "energosbyt/internal/marshaling/config/kafka/producer/broker"
	connection "energosbyt/internal/marshaling/config/kafka/producer/connection"
	retry "energosbyt/internal/marshaling/config/kafka/producer/retry"
)

// Config is the YAML-marshaling representation of Kafka producer configuration.
type Config struct {
	// Brokers stores the marshaling representation of the broker configurations.
	Brokers []broker.Config `yaml:"brokers"`

	// Connection stores the marshaling representation of the connection configuration.
	Connection connection.Config `yaml:"connection"`

	// Retry stores the marshaling representation of the retry configuration.
	Retry retry.Config `yaml:"retry"`

	// Ack stores the required number of acknowledgements from brokers.
	Ack int `yaml:"ack"`
}

// Config converts the marshaling Config into the domain producer.Config.
func (c Config) Config() producer.Config {
	bb := make([]domain.Config, 0, len(c.Brokers))
	for _, b := range c.Brokers {
		bb = append(bb, b.Config())
	}
	return producer.NewConfig(
		bb,
		c.Connection.Config(),
		c.Retry.Config(),
		c.Ack,
	)
}

package kafka

import (
	"energosbyt/internal/config/kafka"
	producer "energosbyt/internal/marshaling/config/kafka/producer"
)

// Config is the YAML-marshaling representation of Kafka configuration.
type Config struct {
	// Producer stores the marshaling representation of the producer configuration.
	Producer producer.Config `yaml:"producer"`
}

// Config converts the marshaling Config into the domain kafka.Config.
func (c Config) Config() kafka.Config {
	return kafka.NewConfig(
		c.Producer.Config(),
	)
}

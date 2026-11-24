package topic

import "energosbyt/internal/config/kafka/producer/broker/topic"

// Config is the YAML-marshaling representation of Kafka topic configuration.
type Config struct {
	// Name stores the Kafka topic name.
	Name string `yaml:"name"`
}

// Config converts the marshaling Config into the domain topic.Config.
func (c Config) Config() topic.Config {
	return topic.NewConfig(
		c.Name,
	)
}

package broker

import (
	"energosbyt/internal/config/kafka/producer/broker"
	domain "energosbyt/internal/config/kafka/producer/broker/topic"
	topic "energosbyt/internal/marshaling/config/kafka/producer/broker/topic"
)

// Config is the YAML-marshaling representation of Kafka broker configuration.
type Config struct {
	// Name stores the broker hostname or IP address.
	Name string `yaml:"name"`

	// Port stores the broker port number.
	Port int `yaml:"port"`

	// Topics stores the marshaling representation of the topic configurations.
	Topics []topic.Config `yaml:"topics"`
}

// Config converts the marshaling Config into the domain broker.Config.
func (c Config) Config() broker.Config {
	tt := make([]domain.Config, 0, len(c.Topics))
	for _, t := range c.Topics {
		tt = append(tt, t.Config())
	}
	return broker.NewConfig(
		c.Name,
		c.Port,
		tt,
	)
}

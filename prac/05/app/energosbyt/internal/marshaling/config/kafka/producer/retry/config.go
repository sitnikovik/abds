package retry

import "energosbyt/internal/config/kafka/producer/retry"

// Config is the YAML-marshaling representation of Kafka producer retry configuration.
type Config struct {
	// Max stores the maximum number of retry attempts.
	Max int `yaml:"max"`
}

// Config converts the marshaling Config into the domain retry.Config.
func (c Config) Config() retry.Config {
	return retry.NewConfig(
		c.Max,
	)
}

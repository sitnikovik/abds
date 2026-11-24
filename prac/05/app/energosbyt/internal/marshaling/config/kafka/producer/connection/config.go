package connection

import "energosbyt/internal/config/kafka/producer/connection"

// Config is the YAML-marshaling representation of Kafka connection configuration.
type Config struct {
	// Timeout stores the connection timeout in milliseconds.
	Timeout int `yaml:"timeout"`
}

// Config converts the marshaling Config into the domain connection.Config.
func (c Config) Config() connection.Config {
	return connection.NewConfig(
		c.Timeout,
	)
}

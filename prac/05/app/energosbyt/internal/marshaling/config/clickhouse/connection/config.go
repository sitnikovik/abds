package connection

import "energosbyt/internal/config/clickhouse/connection"

// Config is the YAML-marshaling representation of ClickHouse connection configuration.
type Config struct {
	// Host stores the ClickHouse server hostname or IP address.
	Host string `yaml:"host"`

	// Port stores the ClickHouse server port number.
	Port int `yaml:"port"`

	// Database stores the target database name.
	Database string `yaml:"database"`
}

// Config converts the marshaling Config into the domain connection.Config.
func (c Config) Config() connection.Config {
	return connection.NewConfig(
		c.Host,
		c.Port,
		c.Database,
	)
}

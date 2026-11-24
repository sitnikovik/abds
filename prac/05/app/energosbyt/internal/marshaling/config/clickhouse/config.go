package clickhouse

import (
	"energosbyt/internal/config/clickhouse"
	auth "energosbyt/internal/marshaling/config/clickhouse/auth"
	connection "energosbyt/internal/marshaling/config/clickhouse/connection"
)

// Config is the YAML-marshaling representation of ClickHouse configuration.
type Config struct {
	// Connection stores the marshaling representation of the connection configuration.
	Connection connection.Config `yaml:"connection"`

	// Auth stores the marshaling representation of the authentication configuration.
	Auth auth.Config `yaml:"auth"`
}

// Config converts the marshaling Config into the domain clickhouse.Config.
func (c Config) Config() clickhouse.Config {
	return clickhouse.NewConfig(
		c.Connection.Config(),
		c.Auth.Config(),
	)
}

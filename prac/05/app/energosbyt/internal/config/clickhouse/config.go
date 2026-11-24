package clickhouse

import (
	"energosbyt/internal/config/clickhouse/auth"
	"energosbyt/internal/config/clickhouse/connection"
)

// Config holds the complete ClickHouse configuration including connection and authentication settings.
type Config struct {
	// connection stores the connection parameters for ClickHouse.
	connection connection.Config

	// auth stores the authentication credentials for ClickHouse.
	auth auth.Config
}

// NewConfig creates a new Config with the given connection and auth settings.
func NewConfig(
	connection connection.Config,
	auth auth.Config,
) Config {
	return Config{
		connection: connection,
		auth:       auth,
	}
}

// Connection returns the connection configuration.
func (c Config) Connection() connection.Config {
	return c.connection
}

// Auth returns the authentication configuration.
func (c Config) Auth() auth.Config {
	return c.auth
}

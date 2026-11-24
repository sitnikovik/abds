package auth

import "energosbyt/internal/config/clickhouse/auth"

// Config is the YAML-marshaling representation of ClickHouse authentication configuration.
type Config struct {
	// Username stores the login username for ClickHouse.
	Username string `yaml:"username"`

	// Password stores the login password for ClickHouse.
	Password string `yaml:"password"`
}

// Config converts the marshaling Config into the domain auth.Config.
func (c Config) Config() auth.Config {
	return auth.NewConfig(
		c.Username,
		c.Password,
	)
}

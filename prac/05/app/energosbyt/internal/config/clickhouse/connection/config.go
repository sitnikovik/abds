package connection

// Config holds ClickHouse connection parameters.
type Config struct {
	// host stores the ClickHouse server hostname or IP address.
	host string

	// port stores the ClickHouse server port number.
	port int

	// database stores the target database name.
	database string
}

// NewConfig creates a new Config with the given connection parameters.
func NewConfig(
	host string,
	port int,
	database string,
) Config {
	return Config{
		host:     host,
		port:     port,
		database: database,
	}
}

// Host returns the ClickHouse server hostname or IP address.
func (c Config) Host() string {
	return c.host
}

// Port returns the ClickHouse server port number.
func (c Config) Port() int {
	return c.port
}

// Database returns the target database name.
func (c Config) Database() string {
	return c.database
}

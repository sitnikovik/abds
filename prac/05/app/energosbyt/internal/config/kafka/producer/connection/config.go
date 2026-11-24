package connection

// Config holds Kafka connection parameters.
type Config struct {
	// timeout stores the connection timeout in milliseconds.
	timeout int
}

// NewConfig creates a new Config with the given timeout value.
func NewConfig(timeout int) Config {
	return Config{
		timeout: timeout,
	}
}

// Timeout returns the connection timeout in milliseconds.
func (c Config) Timeout() int {
	return c.timeout
}

package retry

// Config holds Kafka producer retry configuration.
type Config struct {
	// max stores the maximum number of retry attempts.
	max int
}

// NewConfig creates a new Config with the given maximum retry attempts.
func NewConfig(max int) Config {
	return Config{
		max: max,
	}
}

// Max returns the maximum number of retry attempts.
func (c Config) Max() int {
	return c.max
}

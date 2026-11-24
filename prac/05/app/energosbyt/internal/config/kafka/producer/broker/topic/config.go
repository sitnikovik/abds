package topic

// Config holds Kafka topic configuration.
type Config struct {
	// name stores the name of the Kafka topic.
	name string
}

// NewConfig creates a new Config with the given topic name.
func NewConfig(name string) Config {
	return Config{
		name: name,
	}
}

// Name returns the Kafka topic name.
func (c Config) Name() string {
	return c.name
}

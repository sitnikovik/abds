package streamer

import (
	"energosbyt/internal/config/clickhouse"
	"energosbyt/internal/config/kafka"
	"energosbyt/internal/config/streamer/dataset"
)

// Config holds the complete streamer application configuration.
type Config struct {
	// datasets stores the dataset configuration for data sources.
	datasets dataset.Config

	// kafka stores the Kafka connection and producer configuration.
	kafka kafka.Config

	// clickhouse stores the ClickHouse database connection configuration.
	clickhouse clickhouse.Config
}

// NewConfig creates a new Config with the given streamer parameters.
func NewConfig(
	datasets dataset.Config,
	kafka kafka.Config,
	clickhouse clickhouse.Config,
) Config {
	return Config{
		datasets:   datasets,
		kafka:      kafka,
		clickhouse: clickhouse,
	}
}

// Datasets returns the dataset configuration.
func (c Config) Datasets() dataset.Config {
	return c.datasets
}

// Kafka returns the Kafka configuration.
func (c Config) Kafka() kafka.Config {
	return c.kafka
}

// Clickhouse returns the ClickHouse configuration.
func (c Config) Clickhouse() clickhouse.Config {
	return c.clickhouse
}

package dataset

import "energosbyt/internal/config/streamer/dataset/realty"

// Config holds the streamer dataset configuration including realty settings.
type Config struct {
	// realty stores the realty dataset configuration.
	realty realty.Config
}

// NewConfig creates a new Config with the given realty configuration.
func NewConfig(realty realty.Config) Config {
	return Config{
		realty: realty,
	}
}

// Realty returns the realty dataset configuration.
func (c Config) Realty() realty.Config {
	return c.realty
}

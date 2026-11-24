package dataset

import (
	"energosbyt/internal/config/streamer/dataset"
	"energosbyt/internal/marshaling/config/dataset/realty"
)

// Config is the YAML-marshaling representation of streamer dataset configuration.
type Config struct {
	// Realty stores the marshaling representation of the realty dataset configuration.
	Realty realty.Config `yaml:"realty"`
}

// Config converts the marshaling Config into the domain dataset.Config.
func (c Config) Config() dataset.Config {
	return dataset.NewConfig(
		c.Realty.Config(),
	)
}

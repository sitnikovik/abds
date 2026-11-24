package realty

import (
	"energosbyt/internal/config/streamer/dataset/realty"
	"energosbyt/internal/marshaling/config/file"
)

// Config is the YAML-marshaling representation of realty dataset configuration.
type Config struct {
	// File stores the marshaling representation of the source file configuration.
	File file.Config `yaml:"file"`
}

// Config converts the marshaling Config into the domain realty.Config.
func (c Config) Config() realty.Config {
	return realty.NewConfig(
		c.File.Config(),
	)
}

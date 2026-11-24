package file

import "energosbyt/internal/config/file"

// Config is the YAML-marshaling representation of file configuration.
type Config struct {
	// Path stores the filesystem path to the data file.
	Path string `yaml:"path"`
}

// Config converts the marshaling Config into the domain file.Config.
func (c Config) Config() file.Config {
	return file.NewConfig(
		c.Path,
	)
}

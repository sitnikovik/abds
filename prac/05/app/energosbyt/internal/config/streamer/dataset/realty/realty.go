package realty

import "energosbyt/internal/config/file"

// Config holds the realty dataset configuration including file settings.
type Config struct {
	// file stores the file configuration for loading realty data.
	file file.Config
}

// NewConfig creates a new Config with the given file configuration.
func NewConfig(file file.Config) Config {
	return Config{
		file: file,
	}
}

// File returns the file configuration for loading realty data.
func (c Config) File() file.Config {
	return c.file
}

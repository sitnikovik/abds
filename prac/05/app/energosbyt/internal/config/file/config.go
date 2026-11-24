package file

// Config holds file path configuration.
type Config struct {
	// path stores the filesystem path to the file.
	path string
}

// NewConfig creates a new Config with the given file path.
func NewConfig(path string) Config {
	return Config{
		path: path,
	}
}

// Path returns the filesystem path to the file.
func (c Config) Path() string {
	return c.path
}

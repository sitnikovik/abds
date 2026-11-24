package streamer

import (
	"fmt"
	"os"
)

// Load reads the YAML configuration file from the given path and unmarshals it into a Config.
//
// Returns an error if the file cannot be read or parsed.
func Load(filepath string) (Config, error) {
	bb, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load file: %w", err)
	}
	return Unmarshal(bb)
}

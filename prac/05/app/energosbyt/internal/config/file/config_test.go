package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/file"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with path", func(t *testing.T) {
		t.Parallel()
		cfg := file.NewConfig("/data/test.csv")
		assert.Equal(t, "/data/test.csv", cfg.Path())
	})

	t.Run("creates config with empty path", func(t *testing.T) {
		t.Parallel()
		cfg := file.NewConfig("")
		assert.Equal(t, "", cfg.Path())
	})

	t.Run("preserves complex path", func(t *testing.T) {
		t.Parallel()
		path := "/var/log/streamer/2024/data.csv"
		cfg := file.NewConfig(path)
		assert.Equal(t, path, cfg.Path())
	})
}

package realty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/file"
	"energosbyt/internal/config/streamer/dataset/realty"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with file", func(t *testing.T) {
		t.Parallel()
		fileCfg := file.NewConfig("/data/realty.csv")
		cfg := realty.NewConfig(fileCfg)
		assert.Equal(t, "/data/realty.csv", cfg.File().Path())
	})

	t.Run("preserves file path", func(t *testing.T) {
		t.Parallel()
		fileCfg := file.NewConfig("/var/data/houses.csv")
		cfg := realty.NewConfig(fileCfg)
		assert.Equal(t, "/var/data/houses.csv", cfg.File().Path())
	})
}

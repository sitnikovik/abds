package dataset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/config/file"
	"energosbyt/internal/config/streamer/dataset"
	"energosbyt/internal/config/streamer/dataset/realty"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("creates config with realty", func(t *testing.T) {
		t.Parallel()
		rltyCfg := realty.NewConfig(file.NewConfig("/data/realty.csv"))
		cfg := dataset.NewConfig(rltyCfg)
		assert.Equal(t, "/data/realty.csv", cfg.Realty().File().Path())
	})
}

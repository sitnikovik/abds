package dataset_test

import (
	"testing"

	marshal "energosbyt/internal/marshaling/config/dataset"
	realty "energosbyt/internal/marshaling/config/dataset/realty"
	file "energosbyt/internal/marshaling/config/file"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Realty: realty.Config{
				File: file.Config{
					Path: "/data/realty.csv",
				},
			},
		}
		_ = cfg.Config()
	})
}

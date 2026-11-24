package file_test

import (
	"testing"

	marshal "energosbyt/internal/marshaling/config/file"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	t.Run("converts to domain config", func(t *testing.T) {
		t.Parallel()
		cfg := marshal.Config{
			Path: "/data/test.csv",
		}
		_ = cfg.Config()
	})
}

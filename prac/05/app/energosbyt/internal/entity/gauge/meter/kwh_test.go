package meter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/gauge/meter"
)

func TestKWPH(t *testing.T) {
	t.Parallel()

	t.Run("creates KWPH from uint32", func(t *testing.T) {
		t.Parallel()
		k := meter.NewKWPH(150)
		assert.Equal(t, uint32(150), k.UInt32())
	})

	t.Run("converts to string", func(t *testing.T) {
		t.Parallel()
		k := meter.NewKWPH(150)
		assert.Equal(t, "150", k.String())
	})

	t.Run("handles zero value", func(t *testing.T) {
		t.Parallel()
		k := meter.NewKWPH(0)
		assert.Equal(t, "0", k.String())
	})
}

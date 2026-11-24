package meter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/gauge/meter"
)

func TestMeters(t *testing.T) {
	t.Parallel()

	t.Run("creates meters with day and night", func(t *testing.T) {
		t.Parallel()
		day := meter.NewKWPH(150)
		night := meter.NewKWPH(80)
		m := meter.NewMeters(day, night)
		assert.Equal(t, day, m.Day())
	})

	t.Run("returns correct night value", func(t *testing.T) {
		t.Parallel()
		night := meter.NewKWPH(80)
		m := meter.NewMeters(meter.NewKWPH(150), night)
		assert.Equal(t, night, m.Night())
	})

	t.Run("handles zero values", func(t *testing.T) {
		t.Parallel()
		m := meter.NewMeters(meter.NewKWPH(0), meter.NewKWPH(0))
		assert.Equal(t, uint32(0), m.Day().UInt32())
	})
}

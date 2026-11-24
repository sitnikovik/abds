package meter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/gauge/meter"
)

func TestApproximateForArea(t *testing.T) {
	t.Parallel()

	t.Run("returns values in range for small area", func(t *testing.T) {
		t.Parallel()
		m := meter.ApproximateForArea(30)
		day := m.Day().UInt32()
		assert.GreaterOrEqual(t, day, uint32(100))
	})

	t.Run("returns values in range for medium area", func(t *testing.T) {
		t.Parallel()
		m := meter.ApproximateForArea(50)
		day := m.Day().UInt32()
		assert.GreaterOrEqual(t, day, uint32(175))
	})

	t.Run("returns values in range for large area", func(t *testing.T) {
		t.Parallel()
		m := meter.ApproximateForArea(85)
		day := m.Day().UInt32()
		assert.GreaterOrEqual(t, day, uint32(250))
	})

	t.Run("returns values in range for extra large area", func(t *testing.T) {
		t.Parallel()
		m := meter.ApproximateForArea(150)
		day := m.Day().UInt32()
		assert.GreaterOrEqual(t, day, uint32(350))
	})
}

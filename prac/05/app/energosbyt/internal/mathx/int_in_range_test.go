package mathx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/mathx"
)

func TestIntInRange(t *testing.T) {
	t.Parallel()

	t.Run("returns min when min equals max", func(t *testing.T) {
		t.Parallel()
		result := mathx.IntInRange(5, 5)
		assert.Equal(t, 5, result)
	})

	t.Run("returns value in range for small range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := mathx.IntInRange(1, 3)
			assert.GreaterOrEqual(t, result, 1)
			assert.LessOrEqual(t, result, 3)
		}
	})

	t.Run("returns value in range for large range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := mathx.IntInRange(100, 200)
			assert.GreaterOrEqual(t, result, 100)
			assert.LessOrEqual(t, result, 200)
		}
	})

	t.Run("handles negative range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := mathx.IntInRange(-10, -5)
			assert.GreaterOrEqual(t, result, -10)
			assert.LessOrEqual(t, result, -5)
		}
	})

	t.Run("handles range crossing zero", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := mathx.IntInRange(-5, 5)
			assert.GreaterOrEqual(t, result, -5)
			assert.LessOrEqual(t, result, 5)
		}
	})

	t.Run("returns min for zero-based range", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			result := mathx.IntInRange(0, 10)
			assert.GreaterOrEqual(t, result, 0)
			assert.LessOrEqual(t, result, 10)
		}
	})
}

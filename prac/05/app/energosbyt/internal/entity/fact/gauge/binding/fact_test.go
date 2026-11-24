package binding_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/fact/gauge/binding"
)

func TestFact(t *testing.T) {
	t.Parallel()

	t.Run("creates fact with all fields", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		f := binding.NewFact(123, 456, now)
		assert.Equal(t, uint64(123), f.GaugeID())
	})

	t.Run("returns correct realty ID", func(t *testing.T) {
		t.Parallel()
		f := binding.NewFact(123, 456, time.Now())
		assert.Equal(t, uint64(456), f.RealtyID())
	})

	t.Run("returns correct created at", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		f := binding.NewFact(123, 456, now)
		assert.Equal(t, now, f.CreatedAt())
	})
}

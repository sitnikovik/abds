package gauge_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/gauge"
)

func TestGauge(t *testing.T) {
	t.Parallel()

	t.Run("creates gauge with all fields", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		g := gauge.NewGauge(123, "Test Gauge", "RU", now)
		assert.Equal(t, uint64(123), g.ID())
	})

	t.Run("returns correct name", func(t *testing.T) {
		t.Parallel()
		g := gauge.NewGauge(123, "Test Gauge", "RU", time.Now())
		assert.Equal(t, "Test Gauge", g.Name())
	})

	t.Run("returns correct country", func(t *testing.T) {
		t.Parallel()
		g := gauge.NewGauge(123, "Test Gauge", "RU", time.Now())
		assert.Equal(t, "RU", g.Country())
	})

	t.Run("returns correct created at", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		g := gauge.NewGauge(123, "Test Gauge", "RU", now)
		assert.Equal(t, now, g.CreatedAt())
	})
}

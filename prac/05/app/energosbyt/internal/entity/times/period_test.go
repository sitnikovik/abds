package times_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/times"
)

func TestPeriod(t *testing.T) {
	t.Parallel()

	t.Run("creates period with from and to", func(t *testing.T) {
		t.Parallel()
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(from, to)
		assert.Equal(t, from, p.From())
	})

	t.Run("returns correct to time", func(t *testing.T) {
		t.Parallel()
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(time.Now(), to)
		assert.Equal(t, to, p.To())
	})

	t.Run("creates period with same time", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		p := times.NewPeriod(now, now)
		assert.Equal(t, now, p.From())
	})

	t.Run("contains time inside period", func(t *testing.T) {
		t.Parallel()
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(from, to)
		assert.True(t, p.Contains(time.Date(2020, 3, 15, 12, 0, 0, 0, time.UTC)))
	})

	t.Run("contains period start", func(t *testing.T) {
		t.Parallel()
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(from, to)
		assert.True(t, p.Contains(from))
	})

	t.Run("does not contain time before period", func(t *testing.T) {
		t.Parallel()
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(from, to)
		assert.False(t, p.Contains(from.Add(-time.Second)))
	})

	t.Run("does not contain period end", func(t *testing.T) {
		t.Parallel()
		from := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		p := times.NewPeriod(from, to)
		assert.False(t, p.Contains(to))
	})
}

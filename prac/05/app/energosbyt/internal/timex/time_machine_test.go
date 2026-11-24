package timex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/timex"
)

func TestNewTimeMachine(t *testing.T) {
	t.Parallel()

	t.Run("creates time machine with start time", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			100,
		)
		now := tm.Now()
		assert.Equal(t, start.Year(), now.Year())
	})

	t.Run("creates time machine with speed factor", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			1000,
		)
		time.Sleep(10 * time.Millisecond)
		now := tm.Now()
		assert.True(t, now.After(start))
	})

	t.Run("starts at specified time", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2024,
			6,
			15,
			12,
			30,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			1,
		)
		now := tm.Now()
		assert.Equal(t, 2024, now.Year())
		assert.Equal(t, time.June, now.Month())
		assert.Equal(t, 15, now.Day())
	})
}

func TestSpan(t *testing.T) {
	t.Parallel()

	t.Run("creates time machine with calculated speed", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		span := 5 * 365 * 24 * time.Hour
		real := 20 * time.Second
		tm := timex.Span(
			start,
			span,
			real,
		)
		now := tm.Now()
		assert.Equal(t, start.Year(), now.Year())
	})

	t.Run("simulates years in seconds", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		span := 365 * 24 * time.Hour
		real := 1 * time.Second
		tm := timex.Span(
			start,
			span,
			real,
		)
		time.Sleep(50 * time.Millisecond)
		now := tm.Now()
		assert.True(t, now.After(start))
	})

	t.Run("handles small span", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		span := 1 * time.Hour
		real := 100 * time.Millisecond
		tm := timex.Span(
			start,
			span,
			real,
		)
		time.Sleep(10 * time.Millisecond)
		now := tm.Now()
		assert.True(t, now.After(start))
	})
}

func TestTimeMachine_Now(t *testing.T) {
	t.Parallel()

	t.Run("returns time after start", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			1000,
		)
		time.Sleep(10 * time.Millisecond)
		now := tm.Now()
		assert.True(t, now.After(start))
	})

	t.Run("progresses faster with higher speed", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tmSlow := timex.NewTimeMachine(
			start,
			10,
		)
		tmFast := timex.NewTimeMachine(
			start,
			1000,
		)
		time.Sleep(10 * time.Millisecond)
		nowSlow := tmSlow.Now()
		nowFast := tmFast.Now()
		assert.True(t, nowFast.Sub(start) > nowSlow.Sub(start))
	})

	t.Run("maintains monotonic progression", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			100,
		)
		now1 := tm.Now()
		time.Sleep(5 * time.Millisecond)
		now2 := tm.Now()
		assert.True(t, now2.After(now1))
	})
}

func TestTimeMachine_Interval(t *testing.T) {
	t.Parallel()

	t.Run("returns one second interval", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			100,
		)
		interval := tm.Interval()
		assert.Equal(t, 1*time.Second, interval)
	})

	t.Run("returns consistent interval", func(t *testing.T) {
		t.Parallel()
		start := time.Date(
			2020,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		tm := timex.NewTimeMachine(
			start,
			1000,
		)
		interval1 := tm.Interval()
		interval2 := tm.Interval()
		assert.Equal(t, interval1, interval2)
	})
}

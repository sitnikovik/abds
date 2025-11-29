package time

import (
	"math/rand"
	"time"
)

// Генерирует случайное время в диапазоне [start, end]
func NewTimeInRange(start, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(rand.Int63n(int64(duration)))
	return start.Add(randomDuration)
}

// Генерирует случайную дату между началом и концом года
func NewDateInYear(year int) time.Time {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	return NewTimeInRange(start, end)
}

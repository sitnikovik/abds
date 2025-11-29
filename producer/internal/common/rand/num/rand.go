package num

import "math/rand"

// NewIntInRange генерирует случайное целое число в диапазоне [min, max].
func NewIntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

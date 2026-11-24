package mathx

import "math/rand"

// IntInRange generates a random integer in the inclusive range [min, max].
func IntInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

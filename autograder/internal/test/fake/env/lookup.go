// Package env provides fake environment lookups for tests.
package env

// NewLookup returns a lookup function that reads values from vars,
// with the same contract as os.LookupEnv.
func NewLookup(vars map[string]string) func(key string) (string, bool) {
	return func(key string) (string, bool) {
		val, ok := vars[key]
		return val, ok
	}
}

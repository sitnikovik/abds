// Package step defines the checks of the autograder as ordered units of work.
package step

import (
	"context"
	"fmt"
	"strconv"
)

// Step is a single check of the autograder.
type Step interface {
	// Name returns the short name of the step.
	Name() string
	// Description returns a brief explanation of what the step checks.
	Description() string
	// Run executes the check and returns an error if the check fails.
	//
	// Values produced by the check are stored in vars for the following steps.
	Run(ctx context.Context, vars Variables) error
}

// Variables holds values that steps share during one run.
type Variables struct {
	// values holds the stored values by their keys.
	values map[string]string
}

// NewVariables creates an empty Variables.
func NewVariables() Variables {
	return Variables{
		values: make(map[string]string),
	}
}

// Set stores value under key, replacing the previous value.
func (v Variables) Set(key, value string) {
	v.values[key] = value
}

// Lookup returns the value stored under key and reports whether it exists.
func (v Variables) Lookup(key string) (string, bool) {
	val, ok := v.values[key]
	return val, ok
}

// Require returns the value stored under key, or an error if no earlier step stored it.
func (v Variables) Require(key string) (string, error) {
	val, ok := v.values[key]
	if !ok {
		return "", fmt.Errorf("variable %s is not set by an earlier step", key)
	}
	return val, nil
}

// RequireUint64 returns the value stored under key as an unsigned integer.
//
// It returns an error if no earlier step stored the key or the value is not an unsigned integer.
func (v Variables) RequireUint64(key string) (uint64, error) {
	val, err := v.Require(key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("variable %s is not an unsigned integer: %w", key, err)
	}
	return n, nil
}

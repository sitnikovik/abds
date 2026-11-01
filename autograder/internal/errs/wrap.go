// Package errs provides helpers for annotating the errors of the autograder checks.
package errs

import "fmt"

// Wrap annotates err with msg, so the failing check is named in the report.
//
// It returns nil when err is nil, so it can wrap the result of any check.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

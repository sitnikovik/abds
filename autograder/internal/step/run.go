package step

import "context"

// Result is the outcome of a single step.
type Result struct {
	// name is the name of the step.
	name string
	// err is the error the step returned, nil if the step passed.
	err error
}

// NewResult creates a Result for the step with the given name and error.
func NewResult(name string, err error) Result {
	return Result{
		name: name,
		err:  err,
	}
}

// Name returns the name of the step.
func (r Result) Name() string {
	return r.name
}

// Err returns the error the step returned, nil if the step passed.
func (r Result) Err() error {
	return r.err
}

// RunAll runs the steps in the given order and returns the result of each one.
//
// A failed step does not stop the following steps, so all failures are reported together.
// All steps share one Variables instance.
func RunAll(ctx context.Context, steps []Step) []Result {
	vars := NewVariables()
	results := make([]Result, 0, len(steps))
	for _, s := range steps {
		results = append(results, NewResult(s.Name(), s.Run(ctx, vars)))
	}
	return results
}

package schema

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Run reads the engines and the columns of the database and checks them against the lab spec.
//
// It returns an error if a query fails or any table or column does not match the spec.
func (s Step) Run(
	ctx context.Context,
	_ step.Variables,
) error {
	engines, err := s.db.Strings(
		ctx,
		enginesQuery,
		2,
		s.database,
	)
	if err != nil {
		return fmt.Errorf("read table engines: %w", err)
	}
	columns, err := s.db.Strings(
		ctx,
		columnsQuery,
		3,
		s.database,
	)
	if err != nil {
		return fmt.Errorf("read columns: %w", err)
	}
	engineOf := make(map[string]string, len(engines))
	for _, row := range engines {
		engineOf[row[0]] = row[1]
	}
	typeOf := make(map[string]string, len(columns))
	for _, row := range columns {
		typeOf[row[0]+"."+row[1]] = row[2]
	}
	checks := make([]error, 0)
	for _, t := range tables {
		checks = append(
			checks,
			errs.Wrap(
				strings.Equal(engineOf[t.name], t.engine),
				"table "+t.name+" engine",
			),
		)
		for _, c := range t.columns {
			checks = append(
				checks,
				errs.Wrap(
					strings.Equal(typeOf[t.name+"."+c.name], c.typ),
					"column "+t.name+"."+c.name+" type",
				),
			)
		}
	}
	return expect.All(checks[0], checks[1:]...)
}

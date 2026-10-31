package fixture

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Run reads the first ids of every source table and stores them in vars.
//
// It returns an error if a query fails or a table has fewer rows than its variables need.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	for _, src := range sources {
		query := fmt.Sprintf(
			"SELECT toString(id) FROM %s.%s ORDER BY id LIMIT %d",
			s.database,
			src.table,
			len(src.keys),
		)
		rows, err := s.db.Strings(
			ctx,
			query,
			1,
		)
		if err != nil {
			return fmt.Errorf("read ids of %s: %w", src.table, err)
		}
		err = expect.All(
			errs.Wrap(
				number.GreaterOrEqual(len(rows), len(src.keys)),
				fmt.Sprintf("rows in %s", src.table),
			),
		)
		if err != nil {
			return err
		}
		for i, key := range src.keys {
			vars.Set(key, rows[i][0])
		}
	}
	return nil
}

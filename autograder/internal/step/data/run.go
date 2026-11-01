package data

import (
	"context"
	"fmt"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"

	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Run counts the rows of the tables and checks them against the minimum volumes of the lab spec.
//
// It returns an error if a query fails or a table has fewer rows than required.
func (s Step) Run(
	ctx context.Context,
	_ step.Variables,
) error {
	realties, err := s.Count(
		ctx,
		countQuery,
		realtyTable,
	)
	if err != nil {
		return err
	}
	gauges, err := s.Count(
		ctx,
		countQuery,
		gaugesTable,
	)
	if err != nil {
		return err
	}
	users, err := s.Count(
		ctx,
		countQuery,
		usersTable,
	)
	if err != nil {
		return err
	}
	bindings, err := s.Count(
		ctx,
		countQuery,
		bindingsTable,
	)
	if err != nil {
		return err
	}
	readings, err := s.Count(
		ctx,
		countFinalQuery,
		readingsTable,
	)
	if err != nil {
		return err
	}
	invoices, err := s.Count(
		ctx,
		countQuery,
		invoicesTable,
	)
	if err != nil {
		return err
	}
	replaced, err := s.db.Count(
		ctx,
		fmt.Sprintf(replacedBindingsQuery, s.database),
	)
	if err != nil {
		return fmt.Errorf("count replaced gauge bindings: %w", err)
	}
	tariffs, err := s.db.Count(
		ctx,
		fmt.Sprintf(tariffsWithVersionsQuery, s.database),
	)
	if err != nil {
		return fmt.Errorf("count tariffs with versions: %w", err)
	}
	return expect.All(
		errs.Wrap(
			number.GreaterOrEqual(realties, minRealties),
			"rows in "+realtyTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(gauges, minGauges),
			"rows in "+gaugesTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(users, minUsers),
			"rows in "+usersTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(bindings, minBindings),
			"rows in "+bindingsTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(readings, minReadings),
			"rows in "+readingsTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(invoices, minInvoices),
			"rows in "+invoicesTable,
		),
		errs.Wrap(
			number.GreaterOrEqual(replaced, minReplacedBindings),
			"gauge bindings after a replacement",
		),
		errs.Wrap(
			number.GreaterOrEqual(tariffs, minTariffsWithVersions),
			"tariffs with at least two versions",
		),
	)
}

// Count runs the counting query of the table in the database.
func (s Step) Count(
	ctx context.Context,
	query,
	table string,
) (uint64, error) {
	n, err := s.db.Count(
		ctx,
		fmt.Sprintf(query, s.database, table),
	)
	if err != nil {
		return 0, fmt.Errorf("count %s: %w", table, err)
	}
	return n, nil
}

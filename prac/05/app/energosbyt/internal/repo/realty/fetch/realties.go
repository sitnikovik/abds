package fetch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"energosbyt/internal/entity/realty"
	repo "energosbyt/internal/repo/realty"
)

// Realties retrieves all realty properties from ClickHouse matching the configured filter and limit.
//
// The method builds a SELECT query with optional WHERE clause and LIMIT, scans the results,
// and returns a slice of realty.Realty entities. Returns an error if the query or scanning fails.
func (f Fetch) Realties(ctx context.Context) ([]realty.Realty, error) {
	q := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s FROM %s`,
		"id",
		"address",
		"square",
		"type",
		"created_at",
		repo.Table,
	)
	if !f.ftr.Empty() {
		q += " " + f.ftr.Where()
	}
	if n := f.limit; n > 0 {
		q += ` LIMIT ` + strconv.Itoa(n)
	}
	rows, err := f.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	gg := make([]realty.Realty, 0, f.limit)
	for rows.Next() {
		var (
			id        uint64
			addr      string
			sqr       uint32
			typ       string
			createdAt time.Time
		)
		err = rows.Scan(&id, &addr, &sqr, &typ, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gg = append(gg, realty.NewRealty(
			id,
			addr,
			typ,
			sqr,
			createdAt,
		))
	}
	return gg, nil
}

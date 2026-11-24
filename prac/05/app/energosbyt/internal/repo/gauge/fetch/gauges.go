package fetch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"energosbyt/internal/entity/gauge"
	repo "energosbyt/internal/repo/gauge"
)

// Gauges retrieves all gauges from ClickHouse matching the configured filter and limit.
//
// The method builds a SELECT query with optional WHERE clause and LIMIT, scans the results,
// and returns a slice of gauge.Gauge entities. Returns an error if the query or scanning fails.
func (f Fetch) Gauges(ctx context.Context) ([]gauge.Gauge, error) {
	q := fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s`,
		"id",
		"name",
		"country",
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
	gg := make([]gauge.Gauge, 0, f.limit)
	for rows.Next() {
		var (
			id        uint64
			name      string
			country   string
			createdAt time.Time
		)
		err = rows.Scan(&id, &name, &country, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		gg = append(gg, gauge.NewGauge(
			id,
			name,
			country,
			createdAt,
		))
	}
	return gg, nil
}

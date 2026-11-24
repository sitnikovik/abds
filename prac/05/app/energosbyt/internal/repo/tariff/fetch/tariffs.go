package fetch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"energosbyt/internal/entity/tariff"
	"energosbyt/internal/entity/tariff/price"
	"energosbyt/internal/entity/times"
	repo "energosbyt/internal/repo/tariff"
)

// Tariffs retrieves all tariffs from ClickHouse matching the configured filter and limit.
//
// The method builds a SELECT query with optional WHERE clause and LIMIT, scans the results,
// and returns a slice of tariff.Tariff entities with their pricing and validity periods.
// Returns an error if the query or scanning fails.
func (f Fetch) Tariffs(ctx context.Context) ([]tariff.Tariff, error) {
	q := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s, %s, %s FROM %s`,
		"id",
		"name",
		"price_day_kwph",
		"price_night_kwph",
		"valid_from",
		"valid_to",
		"created_at",
		"updated_at",
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
	trff := make([]tariff.Tariff, 0, f.limit)
	for rows.Next() {
		var (
			id             uint64
			name           string
			priceDayKwph   uint32
			priceNightKwph uint32
			validFrom      time.Time
			validTo        time.Time
			createdAt      time.Time
			updatedAt      time.Time
		)
		err = rows.Scan(
			&id,
			&name,
			&priceDayKwph,
			&priceNightKwph,
			&validFrom,
			&validTo,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		trff = append(trff, tariff.NewTariff(
			id,
			name,
			price.NewPrices(
				priceDayKwph,
				priceNightKwph,
			),
			times.NewPeriod(
				validFrom,
				validTo,
			),
			createdAt,
			updatedAt,
		))
	}
	return trff, nil
}

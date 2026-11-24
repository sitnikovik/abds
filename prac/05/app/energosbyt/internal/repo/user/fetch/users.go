package fetch

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"energosbyt/internal/entity/user"
	repo "energosbyt/internal/repo/user"
)

// Users retrieves all users from ClickHouse matching the configured filter and limit.
//
// The method builds a SELECT query with optional WHERE clause and LIMIT, scans the results,
// and returns a slice of user.User entities. Returns an error if the query or scanning fails.
func (f Fetch) Users(ctx context.Context) ([]user.User, error) {
	q := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s, %s FROM %s`,
		"id",
		"fio",
		"phone",
		"phone_prev",
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
	usrr := make([]user.User, 0, f.limit)
	for rows.Next() {
		var (
			id        uint64
			fio       string
			phone     string
			phonePrev string
			createdAt time.Time
			updatedAt time.Time
		)
		err = rows.Scan(&id, &fio, &phone, &phonePrev, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usrr = append(usrr, user.NewUser(
			id,
			fio,
			phone,
			phonePrev,
			createdAt,
			updatedAt,
		))
	}
	return usrr, nil
}

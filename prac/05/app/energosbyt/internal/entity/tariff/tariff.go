package tariff

import (
	"time"

	"energosbyt/internal/entity/tariff/price"
	"energosbyt/internal/entity/times"
)

// Tariff represents an electricity tariff with its pricing and validity period.
type Tariff struct {
	// valid stores the period during which the tariff is valid.
	valid times.Period

	// createdAt stores the timestamp when the tariff was created.
	createdAt time.Time

	// updatedAt stores the timestamp when the tariff was last updated.
	updatedAt time.Time

	// prices stores the pricing structure for different consumption periods.
	prices price.Prices

	// name stores the display name of the tariff.
	name string

	// id stores the unique identifier of the tariff.
	id uint64
}

// NewTariff creates a new Tariff instance with the given parameters.
func NewTariff(
	id uint64,
	name string,
	prices price.Prices,
	validAt times.Period,
	createdAt time.Time,
	updatedAt time.Time,
) Tariff {
	return Tariff{
		valid:     validAt,
		createdAt: createdAt,
		updatedAt: updatedAt,
		name:      name,
		id:        id,
		prices:    prices,
	}
}

// ID returns the unique identifier of the tariff.
func (t Tariff) ID() uint64 {
	return t.id
}

// IsValidAt reports whether the tariff validity period covers the given moment.
func (t Tariff) IsValidAt(at time.Time) bool {
	return t.valid.Contains(at)
}

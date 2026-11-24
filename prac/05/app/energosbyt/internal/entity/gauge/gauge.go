package gauge

import "time"

// Gauge represents an electricity meter with its identifying information.
type Gauge struct {
	// createdAt stores the timestamp when the gauge was created.
	createdAt time.Time

	// name stores the display name of the gauge.
	name string

	// country stores the country code where the gauge is located.
	country string

	// id stores the unique identifier of the gauge.
	id uint64
}

// NewGauge creates a new Gauge instance with the given parameters.
func NewGauge(
	id uint64,
	name string,
	country string,
	createdAt time.Time,
) Gauge {
	return Gauge{
		createdAt: createdAt,
		name:      name,
		country:   country,
		id:        id,
	}
}

// CreatedAt returns the creation timestamp of the gauge.
func (g Gauge) CreatedAt() time.Time {
	return g.createdAt
}

// Name returns the display name of the gauge.
func (g Gauge) Name() string {
	return g.name
}

// Country returns the country code of the gauge.
func (g Gauge) Country() string {
	return g.country
}

// ID returns the unique identifier of the gauge.
func (g Gauge) ID() uint64 {
	return g.id
}

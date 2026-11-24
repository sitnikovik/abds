package binding

import "time"

// Fact represents a gauge-to-realty binding record.
type Fact struct {
	// createdAt stores the timestamp when the binding was created.
	createdAt time.Time

	// gaugeID stores the identifier of the bound gauge.
	gaugeID uint64

	// realtyID stores the identifier of the bound property.
	realtyID uint64
}

// NewFact creates a new binding Fact with the given parameters.
func NewFact(
	gaugeID uint64,
	realtyID uint64,
	createdAt time.Time,
) Fact {
	return Fact{
		gaugeID:   gaugeID,
		realtyID:  realtyID,
		createdAt: createdAt,
	}
}

// GaugeID returns the identifier of the bound gauge.
func (f Fact) GaugeID() uint64 {
	return f.gaugeID
}

// RealtyID returns the identifier of the bound property.
func (f Fact) RealtyID() uint64 {
	return f.realtyID
}

// CreatedAt returns the creation timestamp of the binding.
func (f Fact) CreatedAt() time.Time {
	return f.createdAt
}

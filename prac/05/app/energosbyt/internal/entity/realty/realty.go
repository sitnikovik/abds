package realty

import "time"

// Realty represents a real estate property in the EnergoSbyt system.
type Realty struct {
	// createdAt stores the timestamp when the property was created.
	createdAt time.Time

	// address stores the street address of the property.
	address string

	// typ stores the type of the property (e.g., apartment, house).
	typ string

	// num stores the apartment or unit number.
	num string

	// id stores the unique identifier of the property.
	id uint64

	// square stores the area of the property in square meters.
	square uint32

	// floor stores the floor number where the property is located.
	floor uint8

	// porch stores the entrance or porch number of the building.
	porch uint8
}

// NewRealty creates a new Realty instance with the given parameters and optional modifiers.
func NewRealty(
	id uint64,
	address string,
	typ string,
	square uint32,
	createdAt time.Time,
	opts ...Option,
) Realty {
	r := Realty{
		id:        id,
		address:   address,
		typ:       typ,
		square:    square,
		createdAt: createdAt,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

// ID returns the unique identifier of the property.
func (r Realty) ID() uint64 {
	return r.id
}

// Address returns the street address of the property.
func (r Realty) Address() string {
	return r.address
}

// Number returns the apartment or unit number.
func (r Realty) Number() string {
	return r.num
}

// Type returns the type of the property.
func (r Realty) Type() string {
	return r.typ
}

// Square returns the area of the property in square meters.
func (r Realty) Square() uint32 {
	return r.square
}

// CreatedAt returns the creation timestamp of the property.
func (r Realty) CreatedAt() time.Time {
	return r.createdAt
}

// Floor returns the floor number where the property is located.
func (r Realty) Floor() uint8 {
	return r.floor
}

// Porch returns the entrance or porch number of the building.
func (r Realty) Porch() uint8 {
	return r.porch
}

package realty

// Option is a functional option for modifying a Realty instance during construction.
type Option func(r *Realty)

// WithNumber returns an Option that sets the apartment or unit number.
func WithNumber(num string) Option {
	return func(r *Realty) {
		r.num = num
	}
}

// WithFloor returns an Option that sets the floor number.
func WithFloor(num uint8) Option {
	return func(r *Realty) {
		r.floor = num
	}
}

// WithPorch returns an Option that sets the entrance or porch number.
func WithPorch(num uint8) Option {
	return func(r *Realty) {
		r.porch = num
	}
}

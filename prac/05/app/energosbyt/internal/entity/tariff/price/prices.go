package price

// Prices holds the day and night electricity tariff prices.
type Prices struct {
	// day stores the daytime electricity price per kilowatt-hour.
	day uint32

	// night stores the nighttime electricity price per kilowatt-hour.
	night uint32
}

// NewPrices creates a new Prices instance with the given day and night rates.
func NewPrices(day, night uint32) Prices {
	return Prices{
		day:   day,
		night: night,
	}
}

// Day returns the daytime electricity price.
func (p Prices) Day() uint32 {
	return p.day
}

// Night returns the nighttime electricity price.
func (p Prices) Night() uint32 {
	return p.night
}

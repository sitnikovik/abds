package meter

// Meters holds day and night electricity meter readings.
type Meters struct {
	// day stores the daytime electricity consumption in kilowatt-hours.
	day KWPH

	// night stores the nighttime electricity consumption in kilowatt-hours.
	night KWPH
}

// NewMeters creates a new Meters instance with the given day and night readings.
func NewMeters(day, night KWPH) Meters {
	return Meters{
		day:   day,
		night: night,
	}
}

// Day returns the daytime electricity reading.
func (m Meters) Day() KWPH {
	return m.day
}

// Night returns the nighttime electricity reading.
func (m Meters) Night() KWPH {
	return m.night
}

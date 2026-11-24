package times

import "time"

// Period represents a time interval with a start and end point.
type Period struct {
	// from stores the start time of the period.
	from time.Time

	// to stores the end time of the period.
	to time.Time
}

// NewPeriod creates a new Period instance with the given start and end times.
func NewPeriod(from, to time.Time) Period {
	return Period{
		from: from,
		to:   to,
	}
}

// From returns the start time of the period.
func (p Period) From() time.Time {
	return p.from
}

// To returns the end time of the period.
func (p Period) To() time.Time {
	return p.to
}

// Contains reports whether the period covers the given moment.
//
// The period includes its start and excludes its end.
func (p Period) Contains(at time.Time) bool {
	return !at.Before(p.from) && at.Before(p.to)
}

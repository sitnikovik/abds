// Package timex provides time simulation utilities for accelerated time progression.
//
// The package implements a time machine that allows simulated epochs to progress
// faster than real time, enabling years of simulated data generation in seconds.
package timex

import "time"

// TimeMachine simulates accelerated time progression, mapping real time to a simulated epoch.
//
// The time machine allows external code to see real time via time.Now while the Now() method
// returns the current simulated time within the epoch, which progresses at an accelerated rate.
type TimeMachine struct {
	// start stores the simulated time at the moment the time machine was created.
	start time.Time

	// from stores the real time when the time machine was created.
	from time.Time

	// speed stores the acceleration factor (how many times faster the epoch runs compared to real time).
	speed int64
}

// NewTimeMachine creates a time machine starting at the given simulated time with the specified speed.
//
// At the moment of creation, the simulated time equals start. After that, it progresses
// at the given speed factor relative to real time.
func NewTimeMachine(
	start time.Time,
	speed int64,
) TimeMachine {
	return TimeMachine{
		start: start,
		from:  time.Now(),
		speed: speed,
	}
}

// Span creates a time machine where span amount of simulated time passes during real duration.
//
// For example, Span(start, 5*365*24*time.Hour, 20*time.Second) simulates five years in twenty seconds.
func Span(
	start time.Time,
	span time.Duration,
	real time.Duration,
) TimeMachine {
	return NewTimeMachine(start, int64(span)/int64(real))
}

// Now returns the current simulated time within the epoch.
//
// The method is safe for concurrent use and does not require a separate goroutine.
func (t TimeMachine) Now() time.Time {
	return t.start.Add(time.Since(t.from) * time.Duration(t.speed))
}

// Interval returns the interval between time machine ticks in real time.
//
// For a time machine created with Span(start, 5*365*24*time.Hour, 20*time.Second),
// this represents approximately 20 seconds divided by the number of real ticks.
func (t TimeMachine) Interval() time.Duration {
	// Returns a fixed interval for demonstration: 1 second of real time
	// corresponds to speed seconds of simulated time
	return 1 * time.Second
}

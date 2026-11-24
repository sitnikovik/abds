package meter

import "energosbyt/internal/mathx"

// ApproximateForArea returns approximate meter readings for the given area in square meters.
//
// The function estimates day and night electricity consumption based on the property area,
// using predefined ranges for different size categories.
func ApproximateForArea(a uint32) Meters {
	var day, night KWPH
	switch {
	case a <= 40:
		day = NewKWPH(uint32(mathx.IntInRange(100, 175)))
		night = NewKWPH(uint32(mathx.IntInRange(45, 75)))
	case a <= 70:
		day = NewKWPH(uint32(mathx.IntInRange(175, 250)))
		night = NewKWPH(uint32(mathx.IntInRange(75, 125)))
	case a <= 100:
		day = NewKWPH(uint32(mathx.IntInRange(250, 350)))
		night = NewKWPH(uint32(mathx.IntInRange(125, 175)))
	default:
		day = NewKWPH(uint32(mathx.IntInRange(350, 500)))
		night = NewKWPH(uint32(mathx.IntInRange(175, 250)))
	}
	return NewMeters(day, night)
}

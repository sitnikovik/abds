package meter

import "strconv"

// KWPH represents kilowatt-hours as an unsigned 32-bit integer.
type KWPH uint32

// NewKWPH creates a new KWPH value from the given uint32.
func NewKWPH(n uint32) KWPH {
	return KWPH(n)
}

// UInt32 returns the KWPH value as a uint32.
func (k KWPH) UInt32() uint32 {
	return uint32(k)
}

// String returns the string representation of the KWPH value.
func (k KWPH) String() string {
	return strconv.FormatUint(uint64(k), 10)
}

package reading

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Fact represents a meter reading record with consumption data.
type Fact struct {
	// createdAt stores the timestamp when the reading was recorded.
	createdAt time.Time

	// gaugeID stores the identifier of the gauge that produced the reading.
	gaugeID uint64

	// realtyID stores the identifier of the property associated with the reading.
	realtyID uint64

	// userID stores the identifier of the user who submitted the reading.
	userID uint64

	// tariffID stores the identifier of the tariff applied to the reading.
	tariffID uint64

	// dayKW stores the daytime electricity consumption in kilowatt-hours.
	dayKW uint32

	// nightKW stores the nighttime electricity consumption in kilowatt-hours.
	nightKW uint32

	// month stores the month for which the reading applies.
	month time.Month
}

// NewFact creates a new meter reading Fact with the given parameters.
func NewFact(
	createdAt time.Time,
	gaugeID uint64,
	realtyID uint64,
	userID uint64,
	tariffID uint64,
	dayKW uint32,
	nightKW uint32,
	month time.Month,
) Fact {
	return Fact{
		createdAt: createdAt,
		gaugeID:   gaugeID,
		realtyID:  realtyID,
		userID:    userID,
		tariffID:  tariffID,
		dayKW:     dayKW,
		nightKW:   nightKW,
		month:     month,
	}
}

// Key returns the gauge identifier as a string key.
func (f Fact) Key() string {
	return strconv.FormatUint(uint64(f.gaugeID), 10)
}

// GaugeID returns the identifier of the gauge.
func (f Fact) GaugeID() uint64 {
	return f.gaugeID
}

// RealtyID returns the identifier of the associated property.
func (f Fact) RealtyID() uint64 {
	return f.realtyID
}

// UserID returns the identifier of the user who submitted the reading.
func (f Fact) UserID() uint64 {
	return f.userID
}

// TariffID returns the identifier of the applied tariff.
func (f Fact) TariffID() uint64 {
	return f.tariffID
}

// DayKW returns the daytime electricity consumption.
func (f Fact) DayKW() uint32 {
	return f.dayKW
}

// NightKW returns the nighttime electricity consumption.
func (f Fact) NightKW() uint32 {
	return f.nightKW
}

// Month returns the month for which the reading applies.
func (f Fact) Month() time.Month {
	return f.month
}

// CreatedAt returns the timestamp when the reading was recorded.
func (f Fact) CreatedAt() time.Time {
	return f.createdAt
}

// MarshalJSON implements json.Marshaler interface.
func (f Fact) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		CreatedAt time.Time `json:"created_at"`
		GaugeID   uint64    `json:"gauge_id"`
		RealtyID  uint64    `json:"realty_id"`
		UserID    uint64    `json:"user_id"`
		TariffID  uint64    `json:"tariff_id"`
		DayKW     uint32    `json:"day_kw"`
		NightKW   uint32    `json:"night_kw"`
		Month     string    `json:"month"`
	}{
		CreatedAt: f.createdAt,
		GaugeID:   f.gaugeID,
		RealtyID:  f.realtyID,
		UserID:    f.userID,
		TariffID:  f.tariffID,
		DayKW:     f.dayKW,
		NightKW:   f.nightKW,
		Month:     f.month.String(),
	})
}

// MustJSON returns the JSON representation of the fact, panicking on error.
func (f Fact) MustJSON() []byte {
	bb, err := f.MarshalJSON()
	if err != nil {
		panic(
			fmt.Sprintf(
				"failed to marshal gauge meter reading fact: %s",
				err.Error(),
			),
		)
	}
	return bb
}

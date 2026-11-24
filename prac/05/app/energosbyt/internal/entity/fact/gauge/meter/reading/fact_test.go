package reading_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/fact/gauge/meter/reading"
)

func TestFact(t *testing.T) {
	t.Parallel()

	t.Run("creates fact with all fields", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		gaugeID := uint64(123)
		realtyID := uint64(456)
		userID := uint64(789)
		tariffID := uint64(10)
		dayKW := uint32(150)
		nightKW := uint32(80)
		month := time.January
		f := reading.NewFact(
			now,
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		assert.Equal(t, gaugeID, f.GaugeID())
	})

	t.Run("returns correct key", func(t *testing.T) {
		t.Parallel()
		gaugeID := uint64(12345)
		realtyID := uint64(67890)
		userID := uint64(11111)
		tariffID := uint64(1)
		dayKW := uint32(100)
		nightKW := uint32(50)
		month := time.February
		f := reading.NewFact(
			time.Now(),
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		assert.Equal(t, "12345", f.Key())
	})

	t.Run("marshals to valid JSON", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2020, 3, 15, 10, 30, 0, 0, time.UTC)
		gaugeID := uint64(123)
		realtyID := uint64(456)
		userID := uint64(789)
		tariffID := uint64(10)
		dayKW := uint32(150)
		nightKW := uint32(80)
		month := time.March
		f := reading.NewFact(
			now,
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		bb := f.MustJSON()
		var result map[string]interface{}
		err := json.Unmarshal(bb, &result)
		assert.NoError(t, err)
	})

	t.Run("JSON contains month as string", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2020, 3, 15, 10, 30, 0, 0, time.UTC)
		gaugeID := uint64(123)
		realtyID := uint64(456)
		userID := uint64(789)
		tariffID := uint64(10)
		dayKW := uint32(150)
		nightKW := uint32(80)
		month := time.March
		f := reading.NewFact(
			now,
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		bb := f.MustJSON()
		var result map[string]interface{}
		_ = json.Unmarshal(bb, &result)
		assert.Equal(t, "March", result["month"])
	})

	t.Run("JSON contains gauge_id", func(t *testing.T) {
		t.Parallel()
		gaugeID := uint64(123)
		realtyID := uint64(456)
		userID := uint64(789)
		tariffID := uint64(10)
		dayKW := uint32(150)
		nightKW := uint32(80)
		month := time.March
		f := reading.NewFact(
			time.Now(),
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		bb := f.MustJSON()
		var result map[string]interface{}
		_ = json.Unmarshal(bb, &result)
		assert.Equal(t, float64(gaugeID), result["gauge_id"])
	})

	t.Run("JSON contains all required fields", func(t *testing.T) {
		t.Parallel()
		gaugeID := uint64(999)
		realtyID := uint64(888)
		userID := uint64(777)
		tariffID := uint64(5)
		dayKW := uint32(200)
		nightKW := uint32(100)
		month := time.July
		f := reading.NewFact(
			time.Now(),
			gaugeID,
			realtyID,
			userID,
			tariffID,
			dayKW,
			nightKW,
			month,
		)
		bb := f.MustJSON()
		var result map[string]interface{}
		_ = json.Unmarshal(bb, &result)
		assert.Contains(t, result, "created_at")
	})
}

package tariff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/tariff"
	"energosbyt/internal/entity/tariff/price"
	"energosbyt/internal/entity/times"
)

func TestTariff(t *testing.T) {
	t.Parallel()

	t.Run("creates tariff with all fields", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		name := "Day-Night"
		prices := price.NewPrices(450, 225)
		validFrom := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		validTo := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		validPeriod := times.NewPeriod(validFrom, validTo)
		now := time.Now()
		updatedAt := now.Add(time.Hour)
		trf := tariff.NewTariff(
			id,
			name,
			prices,
			validPeriod,
			now,
			updatedAt,
		)
		assert.Equal(t, id, trf.ID())
	})

	t.Run("creates single rate tariff", func(t *testing.T) {
		t.Parallel()
		id := uint64(456)
		name := "Single Rate"
		prices := price.NewPrices(500, 0)
		validFrom := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		validTo := time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC)
		validPeriod := times.NewPeriod(validFrom, validTo)
		now := time.Now()
		trf := tariff.NewTariff(
			id,
			name,
			prices,
			validPeriod,
			now,
			now,
		)
		assert.Equal(t, id, trf.ID())
	})
}

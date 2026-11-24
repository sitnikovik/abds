package realty_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/entity/realty"
)

func TestRealty(t *testing.T) {
	t.Parallel()

	t.Run("creates realty with all fields", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		address := "Test Address"
		ty := "apartment"
		square := uint32(50)
		now := time.Now()
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			now,
		)
		assert.Equal(t, id, r.ID())
	})

	t.Run("returns correct address", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		address := "Test Address"
		ty := "apartment"
		square := uint32(50)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
		)
		assert.Equal(t, address, r.Address())
	})

	t.Run("returns correct type", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		address := "Test Address"
		ty := "apartment"
		square := uint32(50)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
		)
		assert.Equal(t, ty, r.Type())
	})

	t.Run("returns correct square", func(t *testing.T) {
		t.Parallel()
		id := uint64(123)
		address := "Test Address"
		ty := "apartment"
		square := uint32(50)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
		)
		assert.Equal(t, square, r.Square())
	})

	t.Run("creates realty with floor option", func(t *testing.T) {
		t.Parallel()
		id := uint64(456)
		address := "Main St"
		ty := "house"
		square := uint32(120)
		floor := uint8(3)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
			realty.WithFloor(floor),
		)
		assert.Equal(t, floor, r.Floor())
	})

	t.Run("creates realty with porch option", func(t *testing.T) {
		t.Parallel()
		id := uint64(456)
		address := "Main St"
		ty := "house"
		square := uint32(120)
		porch := uint8(2)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
			realty.WithPorch(porch),
		)
		assert.Equal(t, porch, r.Porch())
	})

	t.Run("creates realty with number option", func(t *testing.T) {
		t.Parallel()
		id := uint64(456)
		address := "Main St"
		ty := "house"
		square := uint32(120)
		number := "42"
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
			realty.WithNumber(number),
		)
		assert.Equal(t, number, r.Number())
	})

	t.Run("uses default values for optional fields", func(t *testing.T) {
		t.Parallel()
		id := uint64(789)
		address := "Side St"
		ty := "flat"
		square := uint32(30)
		r := realty.NewRealty(
			id,
			address,
			ty,
			square,
			time.Now(),
		)
		assert.Equal(t, uint8(0), r.Floor())
	})
}

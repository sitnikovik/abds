package gauge

import (
	"context"
	"time"

	numRandom "abds-producer/internal/common/rand/num"
	timeRandom "abds-producer/internal/common/rand/time"
	"abds-producer/internal/domain/energosbyt/flat"
	"abds-producer/internal/domain/energosbyt/gauge"
)

// Repo описывает репозиторий для работы с счетчиками.
type Repo struct {
	// m - мапа для хранения последних показаний по счетчику для инкремента.
	m map[gauge.ID]gauge.Gauge
}

// NewRepo создает новый репозиторий для работы с счетчиками.
func NewRepo() *Repo {
	return &Repo{
		m: make(map[gauge.ID]gauge.Gauge),
	}
}

// AllForFlats получает счетчики для квартир.
func (r *Repo) AllForFlats(ctx context.Context, flats flat.Flats) ([]gauge.Gauge, error) {
	n := len(flats)
	res := make([]gauge.Gauge, n)
	for i, fl := range flats {
		id := gauge.NewID(uint32(fl.ID) + 1000)
		last, err := r.LastFor(id)
		if err != nil {
			return nil, err
		}
		var tim int64
		vals := gauge.ApproximateValuesForArea(fl.Area)
		if last.SentAt == 0 {
			tim = timeRandom.NewDateInYear(numRandom.NewIntInRange(2020, 2025)).Unix()
		} else {
			lastTime := time.Unix(last.SentAt, 0)
			tim = timeRandom.NewTimeInRange(
				lastTime.Add(15*24*time.Hour),
				lastTime.Add(20*24*time.Hour),
			).Unix()
			vals.T1 = last.T1.Add(vals.T1)
			vals.T2 = last.T2.Add(vals.T2)
		}
		g := gauge.Gauge{
			SentAt: tim,
			T1:     vals.T1,
			T2:     vals.T2,
			FlatID: fl.ID,
			ID:     id,
		}
		r.Update(ctx, g)
		res[i] = g
	}
	return res, nil
}

// LastFor возвращает текущие показания для счетчика по его идентификатору.
func (r *Repo) LastFor(id gauge.ID) (gauge.Gauge, error) {
	return r.m[id], nil
}

// Update обновляет переданный счетчик.
func (r *Repo) Update(
	_ context.Context,
	g gauge.Gauge,
) error {
	r.m[g.ID] = g
	return nil
}

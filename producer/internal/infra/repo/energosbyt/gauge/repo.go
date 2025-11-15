package gauge

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"abds-producer/internal/domain/energosbyt/gauge"
)

// Repo описывает репозиторий для работы с счетчиками.
type Repo struct{}

// NewRepo создает новый репозиторий для работы с счетчиками.
func NewRepo() *Repo {
	return &Repo{}
}

// All получает все счетчики.
func (r *Repo) All(_ context.Context) ([]gauge.Gauge, error) {
	n := rand.IntN(10) + rand.IntN(10)
	gg := make([]gauge.Gauge, n)
	for i := range n {
		gg[i] = gauge.Gauge{
			T1:     rand.Uint32(),
			T2:     rand.Uint32(),
			FlatID: uint32(i + 1),
			ID:     uint32(i + 1),
		}
	}
	return gg, nil
}

// ListenFor слушает новые счетчики и отправляет их в канал
// пока контекст не будет отменен.
//
// max - буфер канала.
func (r *Repo) ListenFor(
	ctx context.Context,
	max int,
) (<-chan (gauge.Gauge), error) {
	ch := make(chan (gauge.Gauge), max)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				gg, err := r.All(ctx)
				if err != nil {
					continue
				}
				log.Printf("got %d gauges", len(gg))
				for _, g := range gg {
					select {
					case ch <- g:
					case <-ctx.Done():
						return
					}
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(rand.IntN(5)) * time.Second):
				}
			}
		}
	}()
	return ch, nil
}

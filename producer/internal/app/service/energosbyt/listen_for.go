package energosbyt

import (
	"context"
	"log"
	"time"

	"abds-producer/internal/domain/energosbyt/gauge"
)

// ListenFor слушает новые счетчики и отправляет их в канал
// пока контекст не будет отменен.
//
// max - буфер канала.
func (s *Service) ListenFor(
	ctx context.Context,
	max int,
	interval time.Duration,
) (<-chan (gauge.Gauge), error) {
	ch := make(chan (gauge.Gauge), max)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				gg, err := s.AllCurrentGauges(ctx)
				if err != nil {
					log.Printf("failed to get current gauges: %v", err)
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
				case <-time.After(interval):
				}
			}
		}
	}()
	return ch, nil
}

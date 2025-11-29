package listen

import (
	"context"
	"fmt"
	"log"
	"time"

	gaugeMapped "abds-producer/internal/infra/mapper/kafka/gauge"
)

// Listen слушает новые счетчики и отправляет их показания.
//
// max - буфер для хранения счетчиков.
func (u *UseCase) Listen(
	ctx context.Context,
	max int,
	interval time.Duration,
) error {
	gg, err := u.gauges.ListenFor(ctx, max, interval)
	if err != nil {
		return fmt.Errorf("failed to get gauges: %w", err)
	}
	for g := range gg {
		err = u.prod.SendMessage(
			ctx,
			gaugeMapped.ToMessage(u.topic, g),
		)
		if err != nil {
			log.Printf("failed to send message for gauge '%d': %v", g.ID, err)
		} else {
			log.Printf("sent gauge '%d': T1=%d, T2=%d", g.ID, g.T1, g.T2)
		}
	}
	return nil
}

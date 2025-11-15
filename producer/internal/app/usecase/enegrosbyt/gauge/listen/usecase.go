package send

import (
	"context"
	"fmt"
	"log"

	"abds-producer/internal/domain/energosbyt/gauge"
	"abds-producer/internal/infra/client/kafka/message"
	gaugeMapped "abds-producer/internal/infra/mapper/kafka/gauge"
)

// producer реализует продюсера сообщений.
type producer interface {
	// SendMessage отправляет сообщение для консьюмеров.
	SendMessage(
		ctx context.Context,
		msg message.Message,
	) error
}

// gaugeRepo реализует репозиторий счетчиков.
type gaugeRepo interface {
	// ListenFor слушает новые счетчики и отправляет их в канал.
	ListenFor(
		ctx context.Context,
		max int,
	) (<-chan gauge.Gauge, error)
}

// UseCase описывает юзкейc для отправки счетчиков.
type UseCase struct {
	// prod - продюсер сообщений.
	prod producer
	// gauges - репозиторий счетчиков.
	gauges gaugeRepo
	// topic - топик для отправки сообщений.
	topic string
}

// NewUseCase создает новый юзкейc для отправки счетчиков.
//
// prod - продюсер сообщений.
//
// gauges - репозиторий счетчиков.
func NewUseCase(
	prod producer,
	gauges gaugeRepo,
	topic string,
) *UseCase {
	return &UseCase{
		prod:   prod,
		gauges: gauges,
		topic:  topic,
	}
}

// Listen слушает новые счетчики и отправляет их показания.
//
// max - буфер для хранения счетчиков.
func (u *UseCase) Listen(
	ctx context.Context,
	max int,
) error {
	gg, err := u.gauges.ListenFor(ctx, max)
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

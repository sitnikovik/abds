package listen

import (
	"context"
	"time"

	"abds-producer/internal/domain/energosbyt/gauge"
	"abds-producer/internal/infra/client/kafka/message"
)

// producer реализует продюсера сообщений.
type producer interface {
	// SendMessage отправляет сообщение для консьюмеров.
	SendMessage(
		ctx context.Context,
		msg message.Message,
	) error
}

// gaugeListener реализует сервис для работы с счетчиками электроэнергии.
type gaugeListener interface {
	// ListenFor слушает новые счетчики и отправляет их в канал.
	ListenFor(
		ctx context.Context,
		max int,
		interval time.Duration,
	) (<-chan gauge.Gauge, error)
}

// UseCase описывает юзкейc для отправки счетчиков.
type UseCase struct {
	// prod - продюсер сообщений.
	prod producer
	// gauges - сервис счетчиков.
	gauges gaugeListener
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
	gauges gaugeListener,
	topic string,
) *UseCase {
	return &UseCase{
		prod:   prod,
		gauges: gauges,
		topic:  topic,
	}
}

package gauge

import (
	"strconv"

	"abds-producer/internal/domain/energosbyt/gauge"
	"abds-producer/internal/infra/client/kafka/message"
)

// ToMessage мапит счетчик в сообщение для Kafka.
//
// topic - топик, в который будет отправлено сообщение.
//
// g - счетчик электроэнергии, который нужно замапить.
func ToMessage(
	topic string,
	g gauge.Gauge,
) message.Message {
	return message.NewMessage(
		topic,
		strconv.FormatUint(uint64(g.ID), 10),
		g.Bytes(),
	)
}

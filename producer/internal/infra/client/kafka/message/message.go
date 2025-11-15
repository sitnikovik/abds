package message

// Message описывает сообщение для отправки в Kafka.
type Message struct {
	// payload - тело сообщения.
	payload []byte
	// topic - название топика.
	topic string
	// key - ключ сообщения.
	key string
}

// NewMessage создает новое сообщение для отправки в Kafka.
//
// topic - название топика.
//
// key - ключ сообщения.
//
// payload - тело сообщения.
func NewMessage(
	topic,
	key string,
	payload []byte,
) Message {
	return Message{
		topic:   topic,
		key:     key,
		payload: payload,
	}
}

// Topic возвращает название топика.
func (m Message) Topic() string {
	return m.topic
}

// Key возвращает ключ сообщения.
func (m Message) Key() string {
	return m.key
}

// Payload возвращает тело сообщения.
func (m Message) Payload() []byte {
	return m.payload
}

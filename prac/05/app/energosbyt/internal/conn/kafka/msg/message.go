package msg

// Message represents a Kafka message with topic, key, and payload.
type Message struct {
	// payload stores the message content as a byte slice.
	payload []byte

	// topic stores the Kafka topic name for the message.
	topic string

	// key stores the message key for partitioning.
	key string
}

// NewMessage creates a new Kafka Message with the given topic, key, and payload.
func NewMessage(
	topic string,
	key string,
	payload []byte,
) Message {
	return Message{
		topic:   topic,
		key:     key,
		payload: payload,
	}
}

// Topic returns the Kafka topic name for the message.
func (m Message) Topic() string {
	return m.topic
}

// Key returns the message key used for partitioning.
func (m Message) Key() string {
	return m.key
}

// Payload returns the message content as a byte slice.
func (m Message) Payload() []byte {
	return m.payload
}

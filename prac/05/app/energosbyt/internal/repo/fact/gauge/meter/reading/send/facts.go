package send

import (
	"github.com/IBM/sarama"

	"energosbyt/internal/entity/fact/gauge/meter/reading"
)

// Facts manages the sending of meter reading facts to Kafka.
type Facts struct {
	// conn stores the Kafka synchronous producer connection.
	conn sarama.SyncProducer

	// facts stores the list of meter reading facts to send.
	facts []reading.Fact

	// topic stores the Kafka topic name for sending readings.
	topic string
}

// NewFacts creates a new Facts instance for sending the given meter reading facts.
func NewFacts(
	conn sarama.SyncProducer,
	topic string,
	facts ...reading.Fact,
) Facts {
	return Facts{
		conn:  conn,
		facts: facts,
		topic: topic,
	}
}

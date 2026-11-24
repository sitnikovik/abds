package send

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// Run sends all meter reading facts to Kafka.
//
// The method iterates through each fact, creates a Kafka message with the fact's key and JSON payload,
// and sends it to the configured topic. Returns an error if any message fails to send.
func (f Facts) Run(_ context.Context) error {
	for _, fact := range f.facts {
		_, _, err := f.conn.SendMessage(
			&sarama.ProducerMessage{
				Topic: f.topic,
				Key:   sarama.StringEncoder(fact.Key()),
				Value: sarama.ByteEncoder(fact.MustJSON()),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send message to queue broker: %w", err)
		}
	}
	return nil
}

package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"

	"abds-producer/internal/infra/client/kafka/config/producer"
	"abds-producer/internal/infra/client/kafka/message"
)

// Producer описывает синхронного продюсера Kafka.
type Producer struct {
	// conn - подключение к Kafka.
	conn sarama.SyncProducer
	// cfg - конфигурация продюсера.
	cfg producer.Producer
	// connected указывает, установлено ли соединение.
	connected bool
}

// NewProducer создает нового синхронного продюсера Kafka.
func NewProducer(cfg producer.Producer) *Producer {
	return &Producer{
		cfg: cfg,
	}
}

// Connect устанавливает соединение с брокером.
func (p *Producer) Connect() error {
	if p.connected {
		return nil
	}
	cfg := sarama.NewConfig()
	cfg.Net.DialTimeout = time.Duration(p.cfg.Connection.Timeout) * time.Second
	cfg.Metadata.Timeout = 3 * time.Second
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(p.cfg.Acks)
	cfg.Producer.Retry.Max = int(p.cfg.Retry.Max)
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	prod, err := sarama.NewSyncProducer(
		p.cfg.Brokers.Strings(),
		cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to connect broker: %w", err)
	}
	p.conn = prod
	p.connected = true
	return nil
}

// SendMessage отправляет сообщение в Kafka.
func (p *Producer) SendMessage(
	ctx context.Context,
	msg message.Message,
) error {
	if err := p.Connect(); err != nil {
		return fmt.Errorf("failed to connect producer: %w", err)
	}
	_, _, err := p.conn.SendMessage(&sarama.ProducerMessage{
		Topic: msg.Topic(),
		Key:   sarama.StringEncoder(msg.Key()),
		Value: sarama.ByteEncoder(msg.Payload()),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	log.Printf("message sent to topic '%s' with key '%s'", msg.Topic(), msg.Key())
	return nil
}

// Close закрывает соединение с брокером.
func (p *Producer) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}
	return nil
}

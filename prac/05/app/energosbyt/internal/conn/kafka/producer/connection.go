package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"

	config "energosbyt/internal/config/kafka/producer"
)

// Connection establishes a synchronous Kafka producer connection using the given configuration.
//
// The function configures the producer with timeout, acknowledgement, retry, and partitioning settings
// from the provided configuration. Returns a SyncProducer instance or an error if the connection fails.
func Connection(
	ctx context.Context,
	cfg config.Config,
) (sarama.SyncProducer, error) {
	conf := sarama.NewConfig()
	conf.Net.DialTimeout = time.Duration(cfg.Connection().Timeout()) * time.Second
	conf.Metadata.Timeout = 3 * time.Second
	conf.Producer.RequiredAcks = sarama.RequiredAcks(cfg.Ack())
	conf.Producer.Retry.Max = cfg.Retry().Max()
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.Producer.Partitioner = sarama.NewHashPartitioner
	addrs := make([]string, 0, len(cfg.Brokers()))
	for _, brok := range cfg.Brokers() {
		addrs = append(addrs, brok.String())
	}
	return sarama.NewSyncProducer(addrs, conf)
}

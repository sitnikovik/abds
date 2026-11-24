package send

import (
	chgo "github.com/ClickHouse/clickhouse-go/v2"
	kafka "github.com/IBM/sarama"

	"energosbyt/internal/timex"
)

// Process manages the meter reading generation and sending workflow.
type Process struct {
	// db stores the ClickHouse database connection.
	db chgo.Conn

	// queue stores the Kafka synchronous producer for sending messages.
	queue kafka.SyncProducer

	// tm stores the time machine for simulated time progression.
	tm timex.TimeMachine

	// topic stores the Kafka topic name for sending readings.
	topic string
}

// NewProcess creates a new Process with the given database connection, Kafka producer, topic, and time machine.
func NewProcess(
	db chgo.Conn,
	queue kafka.SyncProducer,
	topic string,
	tm timex.TimeMachine,
) Process {
	return Process{
		db:    db,
		queue: queue,
		topic: topic,
		tm:    tm,
	}
}

package producer

import (
	"abds-producer/internal/infra/client/kafka/config/broker"
	"abds-producer/internal/infra/client/kafka/config/connection"
)

// Acks описывает уровень гарантии доставки сообщений.
type Acks int8

const (
	// NoResponse означает, что продюсер не ожидает подтверждения от брокера.
	NoResponse Acks = 0
	// WaitForLocal означает, что продюсер ожидает подтверждения от одной локальной реплики.
	WaitForLocal Acks = 1
	// WaitForAll означает, что продюсер ожидает подтверждения от всех реплик.
	WaitForAll Acks = -1
)

// Producer описывает параметры продюсера.
type Producer struct {
	// Brokers - параметры брокеров, с которыми работает продюсер.
	Brokers broker.Brokers `yaml:"brokers"`
	// Connection - параметры подключения к продюсеру.
	Connection connection.Connection `yaml:"connection"`
	// Retry - параметры ретраев.
	Retry Retry `yaml:"retry"`
	// Acks - уровень гарантии доставки.
	Acks Acks `yaml:"acks"`
}

// Retry описывает параметры ретраев.
type Retry struct {
	// Max - максимальное кол-во ретраев отправки сообщений.
	Max uint8 `yaml:"max"`
}

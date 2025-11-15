package broker

import (
	"fmt"

	"abds-producer/internal/infra/client/kafka/config/topic"
)

// Brokers описывает список параметров брокеров.
type Brokers []Broker

// Broker опиисывает параметры брокера.
type Broker struct {
	// Topics - список топиков, с которыми работает брокер.
	Topics []topic.Topic `yaml:"topics"`
	// Name - имя хоста брокера.
	Name string `yaml:"name"`
	// Port - номер порта брокера.
	Port uint16 `yaml:"port"`
}

// Strings возвращает список адресов брокеров.
func (bb Brokers) Strings() []string {
	ss := make([]string, len(bb))
	for i, b := range bb {
		ss[i] = b.String()
	}
	return ss
}

// String возвращает строковое представление брокера.
func (b Broker) String() string {
	return fmt.Sprintf("%s:%d", b.Name, b.Port)
}

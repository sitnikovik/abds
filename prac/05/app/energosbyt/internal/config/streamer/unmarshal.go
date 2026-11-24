package streamer

import (
	yaml "gopkg.in/yaml.v3"

	clickhouse "energosbyt/internal/marshaling/config/clickhouse"
	dataset "energosbyt/internal/marshaling/config/dataset"
	kafka "energosbyt/internal/marshaling/config/kafka"
)

// yamlConfig accumulates all transport config parts from YAML for unmarshaling.
type yamlConfig struct {
	Dataset    dataset.Config    `yaml:"dataset"`
	Clickhouse clickhouse.Config `yaml:"clickhouse"`
	Kafka      kafka.Config      `yaml:"kafka"`
}

// Unmarshal unmarshals YAML content and returns a Config instance.
//
// Returns an error if the YAML content cannot be parsed.
func Unmarshal(bb []byte) (Config, error) {
	var out yamlConfig
	err := yaml.Unmarshal(bb, &out)
	if err != nil {
		return Config{}, err
	}
	cfg := NewConfig(
		out.Dataset.Config(),
		out.Kafka.Config(),
		out.Clickhouse.Config(),
	)
	return cfg, nil
}

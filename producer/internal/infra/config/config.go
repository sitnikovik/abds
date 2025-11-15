package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config описывает конфигурацию приложения.
type Config struct {
	// Kafka - конфигурация Kafka.
	Kafka Kafka `json:"kafka,omitempty"`
}

// LoadFrom загружает конфигурацию из файла по заданному пути.
func LoadFrom(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

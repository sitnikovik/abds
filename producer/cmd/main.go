package main

import (
	"context"
	"log"

	gauges "abds-producer/internal/app/usecase/enegrosbyt/gauge/listen"
	"abds-producer/internal/infra/client/kafka/producer/sync"
	"abds-producer/internal/infra/config"
	"abds-producer/internal/infra/repo/energosbyt/gauge"
)

func main() {
	f := parseFlags()
	cfg, err := config.LoadFrom(f.configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ctx := context.Background()
	bufsize := 10
	err = gauges.
		NewUseCase(
			sync.NewProducer(cfg.Kafka.Producer),
			gauge.NewRepo(),
			cfg.Kafka.Producer.Brokers[0].Topics[0].Name, // Пока берем первый топик из конфига.
		).
		Listen(ctx, bufsize)
	if err != nil {
		log.Fatalf("failed to send gauges: %v", err)
	}
}

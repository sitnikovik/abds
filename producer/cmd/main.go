package main

import (
	"context"
	"log"
	"time"

	enegrosbyt "abds-producer/internal/app/service/energosbyt"
	gaugesListener "abds-producer/internal/app/usecase/enegrosbyt/gauge/listen"
	"abds-producer/internal/infra/client/kafka/producer/sync"
	"abds-producer/internal/infra/config"
	"abds-producer/internal/infra/repo/energosbyt/flat"
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
	err = gaugesListener.
		NewUseCase(
			sync.NewProducer(cfg.Kafka.Producer),
			enegrosbyt.NewService(
				gauge.NewRepo(),
				flat.NewRepo(f.flatsCSV),
			),
			cfg.Kafka.Producer.Brokers.FirstTopic().Name, // Пока берем первый топик из конфига.
		).
		Listen(ctx, bufsize, 3*time.Second)
	if err != nil {
		log.Fatalf("failed to send gauges: %v", err)
	}
}

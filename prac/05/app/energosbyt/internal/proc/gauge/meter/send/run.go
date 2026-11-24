package send

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"energosbyt/internal/entity/fact/gauge/meter/reading"
	"energosbyt/internal/entity/gauge/meter"
	"energosbyt/internal/repo/fact/gauge/meter/reading/send"
	gauge "energosbyt/internal/repo/gauge/fetch"
	realty "energosbyt/internal/repo/realty/fetch"
	tariff "energosbyt/internal/repo/tariff/fetch"
	user "energosbyt/internal/repo/user/fetch"
)

// Run executes the meter reading generation and sending loop.
//
// The method fetches gauges, realties, users and tariffs from the database once,
// then periodically generates meter readings for each realty
// with the simulated time taken from the time machine
// and sends them to Kafka.
// The loop continues until the context is cancelled,
// so Kafka keeps receiving messages even if the database goes down.
//
// Returns an error if any operation fails.
func (p Process) Run(ctx context.Context) error {
	gauges, err := gauge.
		NewFetch(
			p.db,
			gauge.WithLimit(1_000_000),
		).
		Gauges(ctx)
	if err != nil {
		return fmt.Errorf("failed to select gauges: %w", err)
	}
	if len(gauges) == 0 {
		return fmt.Errorf("there is no any gauge found")
	}
	realties, err := realty.
		NewFetch(
			p.db,
			realty.WithLimit(1_000_000),
		).
		Realties(ctx)
	if err != nil {
		return fmt.Errorf("failed to select realties: %w", err)
	}
	if len(realties) == 0 {
		return fmt.Errorf("there is no any realty found")
	}
	users, err := user.
		NewFetch(
			p.db,
			user.WithLimit(1_000_000),
		).
		Users(ctx)
	if err != nil {
		return fmt.Errorf("failed to select users: %w", err)
	}
	if len(users) == 0 {
		return fmt.Errorf("there is no any user found")
	}
	trff, err := tariff.
		NewFetch(
			p.db,
			tariff.WithLimit(1_000_000),
		).
		Tariffs(ctx)
	if err != nil {
		return fmt.Errorf("failed to select tariffs: %w", err)
	}
	if len(trff) == 0 {
		return fmt.Errorf("there is no any tariff found")
	}
	interval := 1 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Printf(
		"Streamer started. Sending readings every %v for %d realties\n",
		interval,
		len(realties),
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			now := p.tm.Now()
			trf, err := actualTariff(trff, now)
			if err != nil {
				return fmt.Errorf("failed to select actual tariff: %w", err)
			}
			facts := make([]reading.Fact, 0, len(realties))
			for _, rty := range realties {
				gaug := gauges[rand.IntN(len(gauges))]
				usr := users[rand.IntN(len(users))]
				mtr := meter.ApproximateForArea(rty.Square())
				facts = append(facts, reading.NewFact(
					now,
					gaug.ID(),
					rty.ID(),
					usr.ID(),
					trf.ID(),
					mtr.Day().UInt32(),
					mtr.Night().UInt32(),
					now.Month(),
				))
			}
			err = send.
				NewFacts(
					p.queue,
					p.topic,
					facts...,
				).
				Run(ctx)
			if err != nil {
				return fmt.Errorf("failed to send facts to Kafka: %w", err)
			}
			fmt.Printf(
				"Sent %d readings at %v\n",
				len(facts),
				now.Format(time.RFC3339),
			)
		}
	}
}

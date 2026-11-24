package streamer

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	config "energosbyt/internal/config/streamer"
	"energosbyt/internal/conn/clickhouse"
	kafka "energosbyt/internal/conn/kafka/producer"
	"energosbyt/internal/proc/gauge/meter/send"
	"energosbyt/internal/timex"
)

var (
	// cfgFile stores the path to the configuration file.
	cfgFile string

	// rootCmd is the root command for the streamer application.
	rootCmd = &cobra.Command{
		Use:   "streamer",
		Short: "Потоковый генератор показаний счётчиков",
		Long: `Потоковый генератор показаний счётчиков электроэнергии,
		отправляющий данные в брокер Apache Kafka для последующего чтения ClickHouse`,
		RunE: func(
			cmd *cobra.Command,
			args []string,
		) error {
			ctx := context.Background()
			cfg, err := config.Load(cfgFile)
			if err != nil {
				return fmt.Errorf("failed to load config from file: %w", err)
			}
			chconn, err := clickhouse.Connection(
				ctx,
				cfg.Clickhouse(),
			)
			if err != nil {
				return fmt.Errorf("failed to connect ClickHouse: %w", err)
			}
			kafkaconn, err := kafka.Connection(ctx, cfg.Kafka().Producer())
			if err != nil {
				return fmt.Errorf("failed to connect broker: %w", err)
			}
			err = send.
				NewProcess(
					chconn,
					kafkaconn,
					cfg.
						Kafka().
						Producer().
						Brokers()[0].
						Topics()[0].
						Name(),
					timex.Span(
						time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						5*365*24*time.Hour,
						60*time.Second,
					),
				).
				Run(ctx)
			if err != nil {
				return fmt.Errorf("streamer failed: %w", err)
			}
			return nil
		},
	}
)

// Execute executes the root command and returns any error that occurred.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.
		PersistentFlags().
		StringVar(
			&cfgFile,
			"config",
			"streamer.yml",
			"Path to config file",
		)
}

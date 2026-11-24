package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	config "energosbyt/internal/config/clickhouse"
)

// Connection establishes a connection to ClickHouse using the given configuration.
//
// Returns a driver connection instance or an error if the connection fails.
func Connection(
	ctx context.Context,
	cfg config.Config,
) (driver.Conn, error) {
	return clickhouse.Open(
		&clickhouse.Options{
			Addr: []string{
				fmt.Sprintf(
					"%s:%d",
					cfg.
						Connection().
						Host(),
					cfg.
						Connection().
						Port(),
				),
			},
			Auth: clickhouse.Auth{
				Database: cfg.
					Connection().
					Database(),
				Username: cfg.
					Auth().
					Username(),
				Password: cfg.
					Auth().
					Password(),
			},
		},
	)
}

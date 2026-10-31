package clickhouse

import (
	"fmt"

	chgo "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/sitnikovik/abds/autograder/internal/config"
)

// Open creates a Client connected to the ClickHouse server described by cfg.
//
// The connection is established lazily, so Open does not contact the server.
func Open(cfg config.Config) (Client, error) {
	conn, err := chgo.Open(
		&chgo.Options{
			Addr: []string{
				cfg.ClickHouseAddr(),
			},
			Protocol: chgo.Native,
			Auth: chgo.Auth{
				Database: cfg.ClickHouseDatabase(),
				Username: cfg.ClickHouseUser(),
				Password: cfg.ClickHousePassword(),
			},
		},
	)
	if err != nil {
		return Client{}, fmt.Errorf("open connection: %w", err)
	}
	return NewClient(conn), nil
}

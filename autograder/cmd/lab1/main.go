// Command lab1 runs the autograder checks of lab 1 against the student's application and ClickHouse.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sitnikovik/abds/autograder/internal/autograder/lab1"
	"github.com/sitnikovik/abds/autograder/internal/client/clickhouse"
	"github.com/sitnikovik/abds/autograder/internal/config"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// requestTimeout limits the time of every HTTP request of the checks.
const requestTimeout = 5 * time.Second

// main loads the configuration, runs the checks and exits with 1 if any of them failed.
func main() {
	cfg, err := config.Load(os.LookupEnv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	db, err := clickhouse.Open(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to ClickHouse: %v\n", err)
		os.Exit(1)
	}
	steps := lab1.Steps(
		cfg,
		db,
		&http.Client{Timeout: requestTimeout},
	)
	failed := 0
	ctx := context.Background()
	vars := step.NewVariables()
	for _, st := range steps {
		err = st.Run(ctx, vars)
		if err != nil {
			fmt.Printf("❌ %s: %v\n", st.Name(), err)
			continue
		}
		fmt.Printf("✅ %s\n", st.Name())
	}
	if failed > 0 {
		fmt.Printf("%d of %d checks failed\n", failed, len(steps))
		os.Exit(1)
	}
	fmt.Printf("all %d checks passed\n", len(steps))
}

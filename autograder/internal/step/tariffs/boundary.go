package tariffs

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/boolean"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// dateTimeLayout is the layout of the DateTime values in the API and in the database.
const dateTimeLayout = "2006-01-02 15:04:05"

// dayLayout is the layout of the dates in the requests.
const dayLayout = "2006-01-02"

// clickhouseClient is the part of the ClickHouse client the boundary step uses.
type clickhouseClient interface {
	// Strings runs the query and returns the rows with the given number of string columns.
	Strings(ctx context.Context, query string, columns int, args ...any) ([][]string, error)
}

// Boundary checks that the tariff version changes on the day its second version starts (SCD 2).
type Boundary struct {
	// db reads the versions of the tariff.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
	// cli sends the requests.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewBoundary creates a Boundary step that reads the versions from database through db and sends its requests through cli.
func NewBoundary(db clickhouseClient, database string, cli httpx.Doer, baseURL string) Boundary {
	return Boundary{
		db:       db,
		database: database,
		cli:      cli,
		baseURL:  baseURL,
	}
}

// Name returns the name of the step.
func (s Boundary) Name() string {
	return "tariffs boundary"
}

// Description returns a brief explanation of what the step checks.
func (s Boundary) Description() string {
	return "Checks that the tariff version in force changes on the day its second version starts."
}

// Version is the part of the tariff version returned by GET /tariffs the step checks.
type Version struct {
	// ValidFrom is the start of the version.
	ValidFrom string `json:"valid_from"`
	// ValidTo is the end of the version.
	ValidTo string `json:"valid_to"`
}

// Run reads the start of the second version of the tariff and checks the versions in force on the days around it.
//
// It returns an error if a query fails, a request fails or the response does not match the expectations.
func (s Boundary) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.RequireUint64(fixture.TariffID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf("SELECT toString(valid_from) FROM %s.dim_tariffs FINAL WHERE id = ? ORDER BY valid_from", s.database),
		1,
		id,
	)
	if err != nil {
		return fmt.Errorf("read versions: %w", err)
	}
	err = expect.All(errs.Wrap(number.GreaterOrEqual(len(rows), 2), "tariff must have two versions"))
	if err != nil {
		return err
	}
	start, err := time.Parse(dateTimeLayout, rows[1][0])
	if err != nil {
		return fmt.Errorf("parse valid_from: %w", err)
	}
	before := start.AddDate(0, 0, -1).Format(dayLayout)
	after := start.Format(dayLayout)
	prev, err := s.Version(ctx, before)
	if err != nil {
		return err
	}
	cur, err := s.Version(ctx, after)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			boolean.True(prev.ValidFrom <= before+" 00:00:00" && before+" 00:00:00" < prev.ValidTo),
			"version on the day before the boundary must cover that day",
		),
		errs.Wrap(
			boolean.True(cur.ValidFrom <= after+" 00:00:00" && after+" 00:00:00" < cur.ValidTo),
			"version on the boundary day must cover that day",
		),
		errs.Wrap(
			strings.NotEqual(prev.ValidFrom, cur.ValidFrom),
			"version must change on the boundary day",
		),
	)
}

// Version requests the tariff in force on the day and returns its validity.
func (s Boundary) Version(ctx context.Context, day string) (Version, error) {
	req, err := httpx.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/tariffs?date=%s", s.baseURL, day), nil)
	if err != nil {
		return Version{}, err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return Version{}, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return Version{}, err
	}
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return Version{}, errs.Wrap(err, "unexpected response")
	}
	var got Version
	err = httpx.Decode(body, &got)
	if err != nil {
		return Version{}, err
	}
	return got, nil
}

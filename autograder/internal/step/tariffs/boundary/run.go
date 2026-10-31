package boundary

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

// Run reads the start of the second version of the tariff and checks the versions in force on the days around it.
//
// It returns an error if a query fails, a request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.RequireUint64(fixture.TariffID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf(
			"SELECT toString(valid_from) FROM %s.dim_tariffs FINAL WHERE id = ? ORDER BY valid_from",
			s.database,
		),
		1,
		id,
	)
	if err != nil {
		return fmt.Errorf("read versions: %w", err)
	}
	err = expect.All(
		errs.Wrap(
			number.GreaterOrEqual(len(rows), 2),
			"tariff must have two versions",
		),
	)
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
func (s Step) Version(ctx context.Context, day string) (Version, error) {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/tariffs?date=%s", s.baseURL, day),
		nil,
	)
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

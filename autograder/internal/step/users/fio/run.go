package fio

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run changes the full name of the owner and checks the API and the stored rows.
//
// It returns an error if a query fails, a request fails or the result does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.RequireUint64(fixture.UserID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf("SELECT fio FROM %s.dim_users FINAL WHERE id = ?", s.database),
		1,
		id,
	)
	if err != nil {
		return fmt.Errorf("read owner: %w", err)
	}
	err = expect.All(errs.Wrap(number.Positive(len(rows)), "owner must be stored"))
	if err != nil {
		return err
	}
	old := rows[0][0]
	next := "Автограйдер Тест"
	if old == next {
		next = "Автограйдер Тест 2"
	}
	err = s.Patch(ctx, id, next)
	if err != nil {
		return err
	}
	fio, err := s.Fetch(ctx, id)
	if err != nil {
		return err
	}
	stale, err := s.db.Count(
		ctx,
		fmt.Sprintf("SELECT count() FROM %s.dim_users FINAL WHERE id = ? AND fio = ?", s.database),
		id,
		old,
	)
	if err != nil {
		return fmt.Errorf("count previous name: %w", err)
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(fio, next),
			"field fio",
		),
		errs.Wrap(
			number.Equal(stale, 0),
			"previous name must not be stored",
		),
	)
}

// Patch sends the new full name of the owner and expects 204 No Content.
func (s Step) Patch(ctx context.Context, id uint64, fio string) error {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%d", s.baseURL, id),
		map[string]string{"fio": fio},
	)
	if err != nil {
		return err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return err
	}
	return expect.All(
		httpx.Status(rsp, http.StatusNoContent),
		errs.Wrap(
			strings.Empty(string(body)),
			"body must be empty",
		),
	)
}

// Fetch returns the full name of the owner from GET /users/{id}.
func (s Step) Fetch(ctx context.Context, id uint64) (string, error) {
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/%s", s.baseURL, strconv.FormatUint(id, 10)),
		nil,
	)
	if err != nil {
		return "", err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	body, err := httpx.ReadBody(rsp)
	if err != nil {
		return "", err
	}
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return "", errs.Wrap(err, "unexpected response")
	}
	var got struct {
		FIO string `json:"fio"`
	}
	err = json.Unmarshal(body, &got)
	if err != nil {
		return "", fmt.Errorf("decode body: %w", err)
	}
	return got.FIO, nil
}

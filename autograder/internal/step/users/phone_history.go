package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/number"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// PhoneHistory checks that PATCH /users/{id} keeps the previous phone in phone_prev (SCD 3).
type PhoneHistory struct {
	// db reads the stored owner.
	db clickhouseClient
	// database is the name of the database the student's tables live in.
	database string
	// cli sends the requests.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
}

// NewPhoneHistory creates a PhoneHistory step that reads the owner from database through db and sends its requests through cli.
func NewPhoneHistory(db clickhouseClient, database string, cli httpx.Doer, baseURL string) PhoneHistory {
	return PhoneHistory{
		db:       db,
		database: database,
		cli:      cli,
		baseURL:  baseURL,
	}
}

// Name returns the name of the step.
func (s PhoneHistory) Name() string {
	return "users phone history"
}

// Description returns a brief explanation of what the step checks.
func (s PhoneHistory) Description() string {
	return "Checks that PATCH /users/{id} keeps the previous phone in phone_prev."
}

// Run changes the phone of the owner and checks that the previous phone is kept in phone_prev.
//
// It returns an error if a query fails, a request fails or the response does not match the expectations.
func (s PhoneHistory) Run(ctx context.Context, vars step.Variables) error {
	id, err := vars.RequireUint64(fixture.UserID)
	if err != nil {
		return err
	}
	rows, err := s.db.Strings(
		ctx,
		fmt.Sprintf("SELECT phone FROM %s.dim_users FINAL WHERE id = ?", s.database),
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
	next := "+79999999999"
	if old == next {
		next = "+79999999998"
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/users/%d", s.baseURL, id),
		map[string]string{"phone": next},
	)
	if err != nil {
		return err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	patchBody, err := httpx.ReadBody(rsp)
	_ = rsp.Body.Close()
	if err != nil {
		return err
	}
	err = expect.All(httpx.Status(rsp, http.StatusNoContent), errs.Wrap(strings.Empty(string(patchBody)), "body must be empty"))
	if err != nil {
		return err
	}
	got, err := s.Fetch(ctx, id)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(got.Phone, next),
			"field phone",
		),
		errs.Wrap(
			strings.Equal(got.PhonePrev, old),
			"field phone_prev must keep the previous phone",
		),
	)
}

// Owner is the part of the owner returned by GET /users/{id} the step checks.
type Owner struct {
	// Phone is the current phone of the owner.
	Phone string `json:"phone"`
	// PhonePrev is the phone the owner had before the current one.
	PhonePrev string `json:"phone_prev"`
}

// Fetch returns the owner from GET /users/{id}.
func (s PhoneHistory) Fetch(ctx context.Context, id uint64) (Owner, error) {
	req, err := httpx.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/users/%d", s.baseURL, id), nil)
	if err != nil {
		return Owner{}, err
	}
	rsp, err := s.cli.Do(req)
	if err != nil {
		return Owner{}, fmt.Errorf("send request: %w", err)
	}
	body, err := httpx.ReadBody(rsp)
	_ = rsp.Body.Close()
	if err != nil {
		return Owner{}, err
	}
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return Owner{}, errs.Wrap(err, "unexpected response")
	}
	var got Owner
	err = httpx.Decode(body, &got)
	if err != nil {
		return Owner{}, err
	}
	return got, nil
}

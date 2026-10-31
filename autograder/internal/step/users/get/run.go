package get

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/fixture"
)

// Run requests the owner stored under the user id and checks the response.
//
// It returns an error if the request fails or the response does not match the expectations.
func (s Step) Run(
	ctx context.Context,
	vars step.Variables,
) error {
	id, err := vars.Require(fixture.UserID)
	if err != nil {
		return err
	}
	req, err := httpx.NewRequest(
		ctx,
		http.MethodGet,
		s.baseURL+"/users/"+id,
		nil,
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
	err = expect.All(httpx.Status(rsp, http.StatusOK), httpx.JSONContent(rsp))
	if err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	var got struct {
		FIO   string `json:"fio"`
		Phone string `json:"phone"`
		ID    uint64 `json:"id"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return err
	}
	return expect.All(
		errs.Wrap(
			strings.Equal(strconv.FormatUint(got.ID, 10), id),
			"field id",
		),
		errs.Wrap(
			strings.NotEmpty(got.FIO),
			"field fio",
		),
		errs.Wrap(
			strings.NotEmpty(got.Phone),
			"field phone",
		),
	)
}

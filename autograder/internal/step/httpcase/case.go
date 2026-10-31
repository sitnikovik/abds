// Package httpcase provides the declarative step that sends one request and checks its status and error message.
package httpcase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sitnikovik/buggles/expect"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
	"github.com/sitnikovik/abds/autograder/internal/errs"
	"github.com/sitnikovik/abds/autograder/internal/step"
)

// Step sends one request and checks that the response has the expected status and a JSON error message.
type Step struct {
	// desc is the name and the description of the step.
	desc step.Desc
	// cli sends the request.
	cli httpx.Doer
	// baseURL is the base URL of the application, without a trailing slash.
	baseURL string
	// method is the HTTP method of the request.
	method string
	// path is the path of the request, with one %s verb for every variable in keys.
	path string
	// keys are the variables substituted into path, in order.
	keys []string
	// body is the JSON body of the request, nil for no body.
	body any
	// status is the expected status code of the response.
	status int
}

// Option configures a Step.
type Option func(*Step)

// WithMethod sets the HTTP method of the request.
func WithMethod(method string) Option {
	return func(s *Step) {
		s.method = method
	}
}

// WithPath sets the path of the request.
//
// The path is formatted with the values of the variables named in keys, which an earlier step must store.
func WithPath(path string, keys ...string) Option {
	return func(s *Step) {
		s.path = path
		s.keys = keys
	}
}

// WithBody sets the JSON body of the request.
func WithBody(body any) Option {
	return func(s *Step) {
		s.body = body
	}
}

// WithStatus sets the expected status code of the response.
func WithStatus(status int) Option {
	return func(s *Step) {
		s.status = status
	}
}

// NewStep creates a Step with the given description that sends its request to baseURL through cli.
//
// By default the step sends GET to the base URL and expects 200 OK, the options change that.
func NewStep(
	cli httpx.Doer,
	baseURL string,
	desc step.Desc,
	opts ...Option,
) Step {
	s := Step{
		desc:    desc,
		cli:     cli,
		baseURL: baseURL,
		method:  http.MethodGet,
		status:  http.StatusOK,
	}
	for _, opt := range opts {
		opt(&s)
	}
	return s
}

// Name returns the name of the step.
func (s Step) Name() string {
	return s.desc.Name()
}

// Description returns a brief explanation of what the step checks.
func (s Step) Description() string {
	return s.desc.Description()
}

// Run sends the request and checks the status code, the JSON content type and the error message.
//
// It returns an error if a variable is missing, the request fails or the response does not match the expectations.
func (s Step) Run(ctx context.Context, vars step.Variables) error {
	args := make([]any, 0, len(s.keys))
	for _, key := range s.keys {
		val, err := vars.Require(key)
		if err != nil {
			return err
		}
		args = append(args, val)
	}
	req, err := httpx.NewRequest(
		ctx,
		s.method,
		s.baseURL+fmt.Sprintf(s.path, args...),
		s.body,
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
	var got struct {
		Message string `json:"message"`
	}
	err = httpx.Decode(body, &got)
	if err != nil {
		return errs.Wrap(err, "unexpected response")
	}
	return expect.All(
		httpx.Status(rsp, s.status),
		httpx.JSONContent(rsp),
		errs.Wrap(
			strings.NotEmpty(got.Message),
			"field message",
		),
	)
}

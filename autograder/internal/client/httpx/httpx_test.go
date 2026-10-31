package httpx_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/number"
	checkstr "github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/client/httpx"
)

func TestNewRequest(t *testing.T) {
	t.Parallel()

	t.Run("encodes the body as JSON", func(t *testing.T) {
		t.Parallel()
		req, err := httpx.NewRequest(
			context.Background(),
			http.MethodPatch,
			"http://localhost:8080/users/1",
			map[string]string{"phone": "+79000000000"},
		)
		assert.NoError(
			t,
			expect.All(
				checkerr.NoError(err),
				checkstr.Equal(req.Method, http.MethodPatch),
				checkstr.Equal(req.Header.Get("Content-Type"), "application/json"),
			),
		)
	})

	t.Run("creates a request without a body", func(t *testing.T) {
		t.Parallel()
		req, err := httpx.NewRequest(context.Background(), http.MethodGet, "http://localhost:8080/ping", nil)
		assert.NoError(
			t,
			expect.All(
				checkerr.NoError(err),
				checkstr.Empty(req.Header.Get("Content-Type")),
			),
		)
	})

	t.Run("fails on a body that cannot be encoded", func(t *testing.T) {
		t.Parallel()
		_, err := httpx.NewRequest(context.Background(), http.MethodPost, "http://localhost:8080", make(chan int))
		assert.NoError(
			t,
			checkerr.Contains(err, "encode body"),
		)
	})

	t.Run("fails on an invalid url", func(t *testing.T) {
		t.Parallel()
		_, err := httpx.NewRequest(context.Background(), http.MethodGet, "://bad", nil)
		assert.NoError(
			t,
			checkerr.Contains(err, "build request"),
		)
	})
}

func TestReadBody(t *testing.T) {
	t.Parallel()

	t.Run("returns the whole body", func(t *testing.T) {
		t.Parallel()
		body, err := httpx.ReadBody(&http.Response{Body: io.NopCloser(strings.NewReader(`{"ok":true}`))})
		assert.NoError(
			t,
			expect.All(
				checkerr.NoError(err),
				checkstr.Equal(string(body), `{"ok":true}`),
			),
		)
	})
}

func TestStatus(t *testing.T) {
	t.Parallel()

	t.Run("passes on the expected status", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			httpx.Status(&http.Response{StatusCode: http.StatusCreated}, http.StatusCreated),
		)
	})

	t.Run("fails on another status", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			checkerr.Contains(
				httpx.Status(&http.Response{StatusCode: http.StatusNotFound}, http.StatusOK),
				"status code must be 200",
			),
		)
	})

	t.Run("compares numbers", func(t *testing.T) {
		t.Parallel()
		assert.NoError(
			t,
			number.Equal(http.StatusOK, http.StatusOK),
		)
	})
}

func TestJSONContent(t *testing.T) {
	t.Parallel()

	t.Run("passes on a JSON content type", func(t *testing.T) {
		t.Parallel()
		rsp := &http.Response{Header: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}}
		assert.NoError(
			t,
			httpx.JSONContent(rsp),
		)
	})

	t.Run("fails on another content type", func(t *testing.T) {
		t.Parallel()
		rsp := &http.Response{Header: http.Header{"Content-Type": []string{"text/plain"}}}
		assert.NoError(
			t,
			checkerr.Contains(httpx.JSONContent(rsp), "application/json"),
		)
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()

	t.Run("decodes the JSON body", func(t *testing.T) {
		t.Parallel()
		var dst struct {
			Status string `json:"status"`
		}
		err := httpx.Decode([]byte(`{"status":"ok"}`), &dst)
		assert.NoError(
			t,
			expect.All(
				checkerr.NoError(err),
				checkstr.Equal(dst.Status, "ok"),
			),
		)
	})

	t.Run("fails on invalid JSON", func(t *testing.T) {
		t.Parallel()
		var dst struct{}
		assert.NoError(
			t,
			checkerr.Contains(httpx.Decode([]byte(`{`), &dst), "decode body"),
		)
	})
}

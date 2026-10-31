package schema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sitnikovik/buggles/expect"
	checkerr "github.com/sitnikovik/buggles/expect/errors"
	"github.com/sitnikovik/buggles/expect/strings"

	"github.com/sitnikovik/abds/autograder/internal/step"
	"github.com/sitnikovik/abds/autograder/internal/step/schema"
	fakeclickhouse "github.com/sitnikovik/abds/autograder/internal/test/fake/clickhouse"
	fixtureschema "github.com/sitnikovik/abds/autograder/internal/test/fixture/schema"
)

func TestStep_Run(t *testing.T) {
	t.Parallel()

	t.Run("passes on the expected schema", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("system.tables", fixtureschema.Engines()),
				fakeclickhouse.WithStrings("system.columns", fixtureschema.Columns()),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.NoError(err),
		)
	})
	t.Run("fails when no tables exist", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("system.tables", nil),
				fakeclickhouse.WithStrings("system.columns", fixtureschema.Columns()),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.Contains(err, "table dim_realty engine"),
		)
	})
	t.Run("fails when no columns exist", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("system.tables", fixtureschema.Engines()),
				fakeclickhouse.WithStrings("system.columns", nil),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.Contains(err, "column dim_realty.id type"),
		)
	})
	t.Run("fails when reading table engines fails", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithError(errors.New("connection refused")),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			expect.All(
				checkerr.Error(err),
				checkerr.Contains(err, "read table engines"),
			),
		)
	})
	t.Run("fails when reading columns fails", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(
			fakeclickhouse.NewClient(
				fakeclickhouse.WithStrings("system.tables", fixtureschema.Engines()),
			),
			"energosbyt",
		)
		err := s.Run(context.Background(), step.NewVariables())
		assert.NoError(
			t,
			checkerr.Contains(err, "read columns"),
		)
	})
	t.Run("has a name and a description", func(t *testing.T) {
		t.Parallel()
		s := schema.NewStep(fakeclickhouse.NewClient(), "energosbyt")
		assert.NoError(
			t,
			expect.All(
				strings.Equal(s.Name(), "schema"),
				strings.NotEmpty(s.Description()),
			),
		)
	})
}

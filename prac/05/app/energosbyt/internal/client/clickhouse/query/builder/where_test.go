package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"energosbyt/internal/client/clickhouse/query/builder"
)

func TestNewWhere(t *testing.T) {
	t.Parallel()

	t.Run("creates empty where builder", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.NotNil(t, w)
	})

	t.Run("starts with zero count", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, 0, w.Count())
	})

	t.Run("starts with empty args", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, 0, len(w.Args()))
	})

	t.Run("starts with empty string", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, "", w.String())
	})
}

func TestWhere_Add(t *testing.T) {
	t.Parallel()

	t.Run("adds single condition", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("id", uint64(123))
		assert.Equal(t, 1, w.Count())
	})

	t.Run("builds WHERE clause with single condition", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("name", "test")
		assert.Equal(t, "WHERE name = ?", w.String())
	})

	t.Run("preserves argument value", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("id", uint64(456))
		assert.Equal(t, []any{uint64(456)}, w.Args())
	})

	t.Run("adds multiple conditions with AND", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("country", "RU")
		w.Add("status", "active")
		assert.Equal(t, "WHERE country = ? AND status = ?", w.String())
	})

	t.Run("collects all arguments", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("country", "RU")
		w.Add("status", "active")
		assert.Equal(t, 2, len(w.Args()))
	})

	t.Run("increments count correctly", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("field1", "value1")
		w.Add("field2", "value2")
		w.Add("field3", "value3")
		assert.Equal(t, 3, w.Count())
	})

	t.Run("handles different value types", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("id", uint64(1))
		w.Add("name", "test")
		w.Add("active", true)
		args := w.Args()
		assert.Equal(t, uint64(1), args[0])
		assert.Equal(t, "test", args[1])
		assert.Equal(t, true, args[2])
	})
}

func TestWhere_String(t *testing.T) {
	t.Parallel()

	t.Run("returns empty string when no conditions", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, "", w.String())
	})

	t.Run("returns WHERE clause with one condition", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("id", uint64(1))
		assert.Equal(t, "WHERE id = ?", w.String())
	})

	t.Run("returns WHERE clause with multiple conditions", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("a", 1)
		w.Add("b", 2)
		w.Add("c", 3)
		assert.Equal(t, "WHERE a = ? AND b = ? AND c = ?", w.String())
	})
}

func TestWhere_Args(t *testing.T) {
	t.Parallel()

	t.Run("returns empty slice when no conditions", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, 0, len(w.Args()))
	})

	t.Run("returns slice with one argument", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("id", uint64(999))
		assert.Equal(t, []any{uint64(999)}, w.Args())
	})

	t.Run("returns slice with multiple arguments in order", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("first", "a")
		w.Add("second", "b")
		w.Add("third", "c")
		args := w.Args()
		assert.Equal(t, "a", args[0])
		assert.Equal(t, "b", args[1])
		assert.Equal(t, "c", args[2])
	})
}

func TestWhere_Count(t *testing.T) {
	t.Parallel()

	t.Run("returns zero when no conditions added", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		assert.Equal(t, 0, w.Count())
	})

	t.Run("returns correct count after adding conditions", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("x", 1)
		w.Add("y", 2)
		assert.Equal(t, 2, w.Count())
	})

	t.Run("count matches args length", func(t *testing.T) {
		t.Parallel()
		w := builder.NewWhere()
		w.Add("a", 1)
		w.Add("b", 2)
		w.Add("c", 3)
		assert.Equal(t, len(w.Args()), w.Count())
	})
}

package builder

import "strings"

// Where is a WHERE clause builder for ClickHouse queries.
type Where struct {
	// args holds arguments for the WHERE clause.
	args []any

	// sb builds the WHERE clause as a string.
	sb strings.Builder

	// n stores the number of WHERE conditions added.
	n int
}

// NewWhere returns a new Where instance for building WHERE clauses.
func NewWhere() *Where {
	return &Where{
		args: make([]any, 0),
		sb:   strings.Builder{},
	}
}

// Add adds a new WHERE condition with the given field and value.
func (w *Where) Add(field string, value any) {
	if w.n > 0 {
		w.sb.WriteString(" AND ")
	}
	w.sb.WriteString(field)
	w.sb.WriteString(" = ?")
	w.args = append(w.args, value)
	w.n++
}

// String returns the complete WHERE clause as a string.
//
// Returns an empty string if no conditions have been added.
func (w *Where) String() string {
	if w.sb.Len() == 0 {
		return ""
	}
	return "WHERE " + w.sb.String()
}

// Args returns the list of argument values for the WHERE clause.
func (w *Where) Args() []any {
	return w.args
}

// Count returns the number of WHERE conditions that have been added.
func (w *Where) Count() int {
	return w.n
}

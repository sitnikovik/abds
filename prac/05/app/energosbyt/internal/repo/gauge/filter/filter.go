package filter

// where defines the interface for building WHERE clauses.
type where interface {
	// Add adds a new WHERE condition with the given field and value.
	Add(field string, value any)

	// String returns the complete WHERE clause as a string.
	String() string

	// Args returns the list of argument values for the WHERE clause.
	Args() []any
}

// Filter builds WHERE clause filters for gauge queries.
type Filter struct {
	// where stores the underlying WHERE clause builder.
	where where

	// name stores the gauge name to filter by.
	name string
}

// Empty returns true if no filter conditions have been configured.
func (f *Filter) Empty() bool {
	return f.name == ""
}

// Where returns the WHERE clause string for the filter.
//
// The method adds the name condition to the WHERE builder if a name has been set.
func (f *Filter) Where() string {
	if v := f.name; v != "" {
		f.where.Add("name", v)
	}
	return f.where.String()
}

// Args returns the list of argument values for the WHERE clause.
func (f *Filter) Args() []any {
	vv := f.where.Args()
	if len(vv) != 0 {
		return vv
	}
	vv = make([]any, 0, 1)
	if v := f.name; v != "" {
		vv = append(vv, v)
	}
	return vv
}

package step

// Desc holds the name and the description of a step.
type Desc struct {
	// name is the short name of the step.
	name string
	// description is a brief explanation of what the step checks.
	description string
}

// NewDesc creates a Desc with the given name and description.
func NewDesc(name, description string) Desc {
	return Desc{
		name:        name,
		description: description,
	}
}

// Name returns the short name of the step.
func (d Desc) Name() string {
	return d.name
}

// Description returns a brief explanation of what the step checks.
func (d Desc) Description() string {
	return d.description
}

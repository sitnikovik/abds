package building

// Type описывает тип здания.
type Type string

const (
	// Panel - панельное здание.
	Panel Type = "panel"
	// Brick - кирпичное здание.
	Brick Type = "brick"
	// Monolith - монолитное здание.
	Monolith Type = "monolith"
)

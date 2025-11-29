package gauge

// ID описывает идентификатор счетчика электроэнергии.
type ID uint32

// NewID создает идентификатор счетчика электроэнергии.
func NewID[T uint32 | int](num T) ID {
	return ID(num)
}

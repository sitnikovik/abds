package flat

// ID описывает идентификатор квартиры.
type ID uint32

// NewID создает идентификатор квартиры.
func NewID[T uint32 | int](num T) ID {
	return ID(num)
}

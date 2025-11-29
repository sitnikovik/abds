package gauge

// Values описывает показания счетчика электроэнергии.
type Values struct {
	// T1 - показания за дневное время.
	T1 Value `json:"t1"`
	// T2 - показания за ночное время.
	T2 Value `json:"t2"`
}

// Value описывает значение показания счетчика электроэнергии в кВт/ч.
type Value uint32

// NewValue создает значение показания счетчика электроэнергии.
func NewValue[T uint32 | int](num T) Value {
	return Value(num)
}

// Add прибавляет переданное значение к текущему и возвращает новое значение.
func (v Value) Add(value Value) Value {
	return NewValue(uint32(v + value))
}

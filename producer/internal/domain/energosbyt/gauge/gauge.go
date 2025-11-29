package gauge

import (
	"encoding/json"

	numRandom "abds-producer/internal/common/rand/num"
	"abds-producer/internal/domain/energosbyt/flat"
)

// Gauge описывает счетчик электроэнергии.
type Gauge struct {
	// SentAt - время передачи показаний (Unix timestamp в секундах).
	SentAt int64 `json:"sent_at"`
	// T1 - показания дневного потребления электричества.
	T1 Value `json:"t1"`
	// T2 - показания ночного потребления электричества
	T2 Value `json:"t2"`
	// FlatID - идентификатор квартиры, к которой привязан счетчик.
	FlatID flat.ID `json:"flat_id"`
	// ID - идентификатор счетчика.
	ID ID `json:"id"`
}

// Bytes представляет объект в виде слайса байтов.
func (g Gauge) Bytes() []byte {
	bb, _ := json.Marshal(g)
	return bb
}

// ApproximateValuesForArea возвращает приблизительно возможные показания счетчика
// для указанной площади помещения.
func ApproximateValuesForArea(a uint32) Values {
	var vals Values
	switch {
	case a <= 40:
		vals.T1 = Value(numRandom.NewIntInRange(100, 175))
		vals.T2 = Value(numRandom.NewIntInRange(45, 75))
	case a <= 70:
		vals.T1 = Value(numRandom.NewIntInRange(175, 250))
		vals.T2 = Value(numRandom.NewIntInRange(75, 125))
	case a <= 100:
		vals.T1 = Value(numRandom.NewIntInRange(250, 350))
		vals.T2 = Value(numRandom.NewIntInRange(125, 175))
	default:
		vals.T1 = Value(numRandom.NewIntInRange(350, 500))
		vals.T2 = Value(numRandom.NewIntInRange(175, 250))
	}
	return vals
}

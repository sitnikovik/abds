package gauge

import (
	"encoding/json"
)

// Gauge описывает счетчик электроэнергии.
type Gauge struct {
	// T1 - показания дневного потребления электричества.
	T1 uint32 `json:"t1"`
	// T2 - показания ночного потребления электричества.
	T2 uint32 `json:"t2"`
	// FlatID - идентификатор квартиры, к которой привязан счетчик.
	FlatID uint32 `json:"flat_id"`
	// ID - идентификатор счетчика.
	ID uint32 `json:"id"`
}

// Bytes представляет объект в виде слайса байтов.
func (g *Gauge) Bytes() []byte {
	bb, _ := json.Marshal(g)
	return bb
}

package flat

import "abds-producer/internal/domain/energosbyt/building"

// Flat описывает квартиру, по которой ведется учет потребленной электроэнергии.
type Flat struct {
	// Address - адрес квартиры.
	Address string `json:"address,omitempty"`
	// District - район квартиры.
	District string `json:"district,omitempty"`
	// BuildingType - тип здания.
	BuildingType building.Type `json:"building_type,omitempty"`
	// Area - площадь квартиры.
	Area uint32 `json:"area,omitempty"`
	// ID - идентификатор квартиры.
	ID ID `json:"id"`
}

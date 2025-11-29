package energosbyt

import (
	"context"

	"abds-producer/internal/domain/energosbyt/flat"
	"abds-producer/internal/domain/energosbyt/gauge"
)

// gaugeRepo реализует репозиторий счетчиков электроэнергиями.
type gaugeRepo interface {
	// AllForFlats получает счетчики электроэнергии по квартирам.
	AllForFlats(
		ctx context.Context,
		flats flat.Flats,
	) ([]gauge.Gauge, error)
}

// flatsRepo реализует репозиторий квартир.
type flatsRepo interface {
	// All возвращает список всех доступных квартир.
	All(
		ctx context.Context,
	) (flat.Flats, error)
}

// Service описывает сервис для работы с счетчиками электроэнергии.
type Service struct {
	// gauges - репозиторий счетчиков электроэнергии.
	gauges gaugeRepo
	// flats - репозиторий квартир.
	flats flatsRepo
}

// NewService создает сервис для работы с счетчиками электроэнергии.
func NewService(
	gauges gaugeRepo,
	flats flatsRepo,
) *Service {
	return &Service{
		gauges: gauges,
		flats:  flats,
	}
}

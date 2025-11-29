package energosbyt

import (
	"context"
	"fmt"

	"abds-producer/internal/domain/energosbyt/gauge"
)

// AllCurrentGauges возвращает все счетчики электроэнергии
// с актуальными показаниями по всем квартирам.
func (s *Service) AllCurrentGauges(
	ctx context.Context,
) ([]gauge.Gauge, error) {
	flats, err := s.flats.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get flats: %w", err)
	}
	gauges, err := s.gauges.AllForFlats(ctx, flats)
	if err != nil {
		return nil, fmt.Errorf("failed to get gauges: %w", err)
	}
	return gauges, nil
}

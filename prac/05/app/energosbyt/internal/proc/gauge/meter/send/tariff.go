package send

import (
	"fmt"
	"time"

	tariff "energosbyt/internal/entity/tariff"
)

// actualTariff returns the first tariff whose validity period covers the given moment.
//
// Returns an error if no tariff is valid at the given moment.
func actualTariff(trff []tariff.Tariff, at time.Time) (tariff.Tariff, error) {
	for _, trf := range trff {
		if trf.IsValidAt(at) {
			return trf, nil
		}
	}
	return tariff.Tariff{}, fmt.Errorf("there is no any tariff found")
}

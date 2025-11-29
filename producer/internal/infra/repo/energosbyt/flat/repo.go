package flat

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"abds-producer/internal/domain/energosbyt/building"
	"abds-producer/internal/domain/energosbyt/flat"
)

// Repo описывает репозиторий для работы с квартирами.
type Repo struct {
	// filepath - путь к файлу с данными.
	filepath string
}

// NewRepo создает новый репозиторий для работы с квартирами.
func NewRepo(
	filepath string,
) *Repo {
	return &Repo{
		filepath: filepath,
	}
}

// All возвращает список всех доступных квартир.
func (r *Repo) All(_ context.Context) (flat.Flats, error) {
	file, err := os.Open(r.filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}
	flats := make(flat.Flats, 0, len(records))
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 5 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		area, err := strconv.ParseUint(record[3], 10, 32)
		if err != nil {
			continue
		}
		flats = append(flats, flat.Flat{
			ID:           flat.NewID(id),
			Address:      record[1],
			District:     record[2],
			Area:         uint32(area),
			BuildingType: building.Type(record[4]),
		})
	}
	return flats, nil
}

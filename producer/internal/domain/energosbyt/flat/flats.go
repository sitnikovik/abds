package flat

type Flats []Flat

func (ff Flats) IDs() []ID {
	if ff.Empty() {
		return nil
	}
	res := make([]ID, ff.Count())
	for i, f := range ff {
		res[i] = f.ID
	}
	return res
}

func (ff Flats) Empty() bool {
	return len(ff) == 0
}

func (ff Flats) Count() int {
	return len(ff)
}

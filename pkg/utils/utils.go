package utils

func UintToUint64Slice(ids []uint) (res []uint64) {
	res = make([]uint64, 0, len(ids))
	for _, id := range ids {
		res = append(res, uint64(id))
	}

	return
}

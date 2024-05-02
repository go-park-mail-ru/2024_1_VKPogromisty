package utils

func UintToUint64Slice(ids []uint) (res []uint64) {
	for _, id := range ids {
		res = append(res, uint64(id))
	}

	return
}

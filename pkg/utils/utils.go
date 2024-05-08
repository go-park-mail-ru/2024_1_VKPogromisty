package utils

import "strconv"

func UintToUint64Slice(ids []uint) (res []uint64) {
	res = make([]uint64, 0, len(ids))
	for _, id := range ids {
		res = append(res, uint64(id))
	}

	return
}

func UintArrayIntoString(arr []uint) (res string) {
	res = ""
	for i, id := range arr {
		res += strconv.Itoa(int(id))
		if i != len(arr)-1 {
			res += ", "
		}
	}

	return
}

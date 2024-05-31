package utils

import (
	"strconv"

	"github.com/jackc/pgtype"
)

func Uint64ToUintSlice(ids []uint64) (res []uint) {
	res = make([]uint, 0)

	for _, id := range ids {
		res = append(res, uint(id))
	}

	return
}

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

func TextArrayIntoStringSlice(arr pgtype.TextArray) (res []string) {
	for _, v := range arr.Elements {
		if v.Status == pgtype.Present {
			res = append(res, v.String)
		}
	}

	return
}

func Int8ArrayIntoUintSlice(arr pgtype.Int8Array) (res []uint64) {
	for _, v := range arr.Elements {
		if v.Status == pgtype.Present {
			res = append(res, uint64(v.Int))
		}
	}

	return
}

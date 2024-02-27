package utils

import (
	"math/rand"
	"time"
)

type CustomTime struct {
	time.Time
}

var (
	Runes      = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	DateFormat = "2006-01-02"
)

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = Runes[rand.Intn(len(Runes))]
	}
	return string(b)
}

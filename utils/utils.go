package utils

import (
	"socio/errors"
	"time"
)

type CustomTime struct {
	time.Time
}

var (
	DateFormat = "2006-01-02"
)

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		err = errors.ErrInvalidDate
		return err
	}

	t.Time = date
	return
}

package customtime

import (
	"socio/errors"
	"time"
)

var (
	DateFormat = "2006-01-02"
)

type TimeProvider interface {
	Now() time.Time
}

type MockTimeProvider struct{}

func (MockTimeProvider) Now() time.Time {
	return time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		err = errors.ErrInvalidDate
		return err
	}

	t.Time = date
	return
}

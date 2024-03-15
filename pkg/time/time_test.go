package customtime_test

import (
	"socio/errors"
	customtime "socio/pkg/time"
	"testing"
	"time"
)

type TimeTestCase struct {
	Data         []byte
	Err          error
	ExpectedTime time.Time
}

var TimeTestCases = map[string]TimeTestCase{
	"valid time": {
		Data:         []byte(`"2006-01-02T15:04:05.000-0700"`),
		Err:          nil,
		ExpectedTime: time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
	},
	"invalid time": {
		Data:         []byte(`""`),
		Err:          errors.ErrInvalidDate,
		ExpectedTime: time.Time{},
	},
}

func TestCustomTimeUnmarshalJSON(t *testing.T) {
	for name, tc := range TimeTestCases {
		t.Run(name, func(t *testing.T) {
			var customTime customtime.CustomTime
			err := customTime.UnmarshalJSON(tc.Data)

			if err != tc.Err || customTime.Time != tc.ExpectedTime {
				t.Errorf("wrong customTime: got %s, expected %s", customTime.Time, tc.ExpectedTime)
				return
			}
		})
	}
}

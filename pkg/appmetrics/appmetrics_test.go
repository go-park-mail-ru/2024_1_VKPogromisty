package appmetrics

import (
	"errors"
	"testing"

	customtime "socio/pkg/time"
)

func TestTrackAppExternalServiceMetrics(t *testing.T) {
	tp := customtime.MockTimeProvider{}

	tests := []struct {
		name       string
		systemName string
		startTime  customtime.CustomTime
		err        error
	}{
		{
			name:       "Test with no error",
			systemName: "TestSystem",
			startTime: customtime.CustomTime{
				Time: tp.Now(),
			},
			err: nil,
		},
		{
			name:       "Test with error",
			systemName: "TestSystem",
			startTime: customtime.CustomTime{
				Time: tp.Now(),
			},
			err: errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TrackAppExternalServiceMetrics(tt.systemName, tt.startTime, tt.err)
		})
	}
}

package utils

import (
	"reflect"
	"testing"
)

func TestUintToUint64Slice(t *testing.T) {
	tests := []struct {
		name     string
		ids      []uint
		expected []uint64
	}{
		{
			name:     "Test empty slice",
			ids:      []uint{},
			expected: []uint64{},
		},
		{
			name:     "Test non-empty slice",
			ids:      []uint{1, 2, 3},
			expected: []uint64{1, 2, 3},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := UintToUint64Slice(tt.ids)

			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("UintToUint64Slice(%v) = %v, want %v", tt.ids, res, tt.expected)
			}
		})
	}
}

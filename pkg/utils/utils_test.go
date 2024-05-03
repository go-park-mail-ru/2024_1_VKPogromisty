package utils

import (
	"reflect"
	"testing"
)

func TestUintToUint64Slice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []uint
		expected []uint64
	}{
		{
			name:     "Test empty array",
			input:    []uint{},
			expected: []uint64{},
		},
		{
			name:     "Test single element array",
			input:    []uint{1},
			expected: []uint64{1},
		},
		{
			name:     "Test multiple elements array",
			input:    []uint{1, 2, 3},
			expected: []uint64{1, 2, 3},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := UintToUint64Slice(tt.input)

			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, output)
			}
		})
	}
}

func TestUintArrayIntoString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []uint
		expected string
	}{
		{
			name:     "Test empty array",
			input:    []uint{},
			expected: "",
		},
		{
			name:     "Test single element array",
			input:    []uint{1},
			expected: "1",
		},
		{
			name:     "Test multiple elements array",
			input:    []uint{1, 2, 3},
			expected: "1, 2, 3",
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := UintArrayIntoString(tt.input)

			if output != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, output)
			}
		})
	}
}

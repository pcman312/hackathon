package services

import (
	"testing"
)

func TestAddOne(t *testing.T) {
	tcs := []struct {
		input    int
		expected int
	}{
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 5},
		{5, 6},
		{6, 7},
		{7, 8},
		{8, 9},
		{9, 10},
		{10, 11},
	}

	for _, tt := range tcs {
		actual := addOne(tt.input)
		if actual != tt.expected {
			t.Fatalf("Expected [%d] but was [%d]", tt.expected, actual)
		}
	}
}

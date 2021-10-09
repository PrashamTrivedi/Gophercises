package main

import (
	"fmt"
	"testing"
)

func TestNormalize(t *testing.T) {

	testCase := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, tc := range testCase {
		t.Run(fmt.Sprintf("Normalizing %s", tc.input), func(t *testing.T) {
			actual := normalizePhoneNo(tc.input)
			if actual != tc.want {
				t.Errorf("got %s; want %s", actual, tc.want)
			}

		})
	}
}

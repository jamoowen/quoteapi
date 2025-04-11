package utils

import (
	"testing"
)

func TestLooselyCompareTwoStrings(t *testing.T) {
	// Defining the columns of the table
	var tests = []struct {
		name string
		s1   string
		s2   string
		want bool
	}{
		// the table itself
		{"Similar strings should return true", "James", "jame", true},
		{"Similar strings with spaces should return true", "James Owen", "james owen", true},
		{"Similar strings with extra spaces should return true", "James Owen", "james   owen", true},
		{"Dissimilar strings should return false", "James", "jane", false},
		{"If nothing is supplied should return false", "James", "", false},
		{"If extra input is given in candidate string should return false", "James", "jameso", false},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := LooselyCompareTwoStrings(tt.s1, tt.s2)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

package main

import "testing"

func TestFormatRating(t *testing.T) {
	tests := []struct {
		Rating   float64
		Expected string
	}{
		{
			Rating:   1,
			Expected: "★",
		},
		{
			Rating:   4,
			Expected: "★★★★",
		},
		{
			Rating:   4.5,
			Expected: "★★★★ ½",
		},
	}

	for _, test := range tests {
		got := convertToStars(test.Rating)
		if got != test.Expected {
			t.Errorf("got %s, expected %s for rating %s", got, test.Expected, test.Rating)
		}
	}
}

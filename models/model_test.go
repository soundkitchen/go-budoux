package models

import (
	"fmt"
	"testing"
)

func TestModelTotalScore(t *testing.T) {
	cases := []struct {
		Expected float64
		Model    Model
	}{
		{
			Expected: 10,
			Model: Model{
				"good": {
					"morning":   1,
					"afternoon": 2,
				},
				"bad": {
					"morning":   3,
					"afternoon": 4,
				},
			},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			score := c.Model.TotalScore()
			if score != c.Expected {
				t.Errorf("Expected %f, got %f", c.Expected, score)
			}
		})
	}
}

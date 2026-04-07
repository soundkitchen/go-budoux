// Copyright 2021-2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"fmt"
	"testing"
)

func TestModelClone(t *testing.T) {
	original := Model{
		"group": {
			"feature": 1,
		},
	}

	cloned := original.Clone()
	cloned["group"]["feature"] = 2

	if original["group"]["feature"] != 1 {
		t.Fatalf("Clone should deep copy nested maps; original was mutated to %d", original["group"]["feature"])
	}
}

func TestModelTotalScore(t *testing.T) {
	cases := []struct {
		Expected int
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
				t.Errorf("Expected %d, got %d", c.Expected, score)
			}
		})
	}
}

func TestGetDefaultJapaneseModelReturnsClone(t *testing.T) {
	modelA := GetDefaultJapaneseModel()

	var (
		groupKey   string
		featureKey string
		original   int
		found      bool
	)
	for group, features := range modelA {
		for feature, score := range features {
			groupKey = group
			featureKey = feature
			original = score
			found = true
			break
		}
		if found {
			break
		}
	}
	if !found {
		t.Fatal("default japanese model should not be empty")
	}

	modelA[groupKey][featureKey] = original + 1
	modelB := GetDefaultJapaneseModel()
	if modelB[groupKey][featureKey] != original {
		t.Fatalf("default model should be cloned; expected %d, got %d", original, modelB[groupKey][featureKey])
	}
}

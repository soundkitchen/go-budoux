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

import "maps"

type Model map[string]map[string]int

// Clone deep-copies the model and all feature groups.
func (m Model) Clone() Model {
	if m == nil {
		return nil
	}

	cloned := make(Model, len(m))
	for group, features := range m {
		groupClone := make(map[string]int, len(features))
		for feature, score := range features {
			groupClone[feature] = score
		}
		cloned[group] = groupClone
	}
	return cloned
}

// TotalScore calculates the total score of the model.
func (m Model) TotalScore() int {
	var total int
	for group := range maps.Values(m) {
		for score := range maps.Values(group) {
			total += score
		}
	}
	return total
}

func GetDefaultJapaneseModel() Model {
	return ja.Clone()
}

func GetJapaneseKNBCModel() Model {
	return ja_knbc.Clone()
}

func GetDefaultThaiModel() Model {
	return th.Clone()
}

func GetDefaultSimplifiedChineseModel() Model {
	return zhhans.Clone()
}

func GetDefaultTraditionalChineseModel() Model {
	return zhhant.Clone()
}

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

package budoux

import (
	"unicode/utf8"

	"github.com/soundkitchen/go-budoux/models"
)

type Parser struct {
	model     models.Model
	baseScore float64

	uw1 map[string]float64
	uw2 map[string]float64
	uw3 map[string]float64
	uw4 map[string]float64
	uw5 map[string]float64
	uw6 map[string]float64

	bw1 map[string]float64
	bw2 map[string]float64
	bw3 map[string]float64

	tw1 map[string]float64
	tw2 map[string]float64
	tw3 map[string]float64
	tw4 map[string]float64
}

// Create new parser.
func New(model models.Model) *Parser {
	return &Parser{
		model:     model,
		baseScore: -model.TotalScore() * 0.5,
		uw1:       model["UW1"],
		uw2:       model["UW2"],
		uw3:       model["UW3"],
		uw4:       model["UW4"],
		uw5:       model["UW5"],
		uw6:       model["UW6"],
		bw1:       model["BW1"],
		bw2:       model["BW2"],
		bw3:       model["BW3"],
		tw1:       model["TW1"],
		tw2:       model["TW2"],
		tw3:       model["TW3"],
		tw4:       model["TW4"],
	}
}

// NewDefaultJapaneseParser returns new Parser with default japanese model.
func NewDefaultJapaneseParser() *Parser {
	return New(models.GetDefaultJapaneseModel())
}

// NewJapaneseKNBCParser returns new Parser with the japanese KNBC base model.
func NewJapaneseKNBCParser() *Parser {
	return New(models.GetJapaneseKNBCModel())
}

// NewDefaultThaiParser returns new Parser with default thai model.
func NewDefaultThaiParser() *Parser {
	return New(models.GetDefaultThaiModel())
}

// NewDefaultSimplifiedChineseParser returns new Parser with default simplified chinese model.
func NewDefaultSimplifiedChineseParser() *Parser {
	return New(models.GetDefaultSimplifiedChineseModel())
}

// NewDefaultTraditionalChineseParser returns new Parser with default traditional chinese model.
func NewDefaultTraditionalChineseParser() *Parser {
	return New(models.GetDefaultTraditionalChineseModel())
}

func lookupScore(group map[string]float64, key string) float64 {
	if group == nil {
		return 0
	}
	return group[key]
}

// Parses a sentence into phrases.
func (p *Parser) Parse(sentence string) []string {
	if sentence == "" {
		return []string{}
	}

	offsets := make([]int, 0, utf8.RuneCountInString(sentence)+1)
	for i := range sentence {
		offsets = append(offsets, i)
	}
	offsets = append(offsets, len(sentence))

	if len(offsets) == 2 {
		return []string{sentence}
	}

	boundaries := make([]int, 1, len(offsets))

	for i := 1; i < len(offsets)-1; i++ {
		score := p.baseScore
		if i > 2 {
			score += lookupScore(p.uw1, sentence[offsets[i-3]:offsets[i-2]])
		}
		if i > 1 {
			score += lookupScore(p.uw2, sentence[offsets[i-2]:offsets[i-1]])
		}
		score += lookupScore(p.uw3, sentence[offsets[i-1]:offsets[i]])
		score += lookupScore(p.uw4, sentence[offsets[i]:offsets[i+1]])
		if i+1 < len(offsets)-1 {
			score += lookupScore(p.uw5, sentence[offsets[i+1]:offsets[i+2]])
		}
		if i+2 < len(offsets)-1 {
			score += lookupScore(p.uw6, sentence[offsets[i+2]:offsets[i+3]])
		}

		if i > 1 {
			score += lookupScore(p.bw1, sentence[offsets[i-2]:offsets[i]])
		}
		score += lookupScore(p.bw2, sentence[offsets[i-1]:offsets[i+1]])
		if i+1 < len(offsets)-1 {
			score += lookupScore(p.bw3, sentence[offsets[i]:offsets[i+2]])
		}

		if i > 2 {
			score += lookupScore(p.tw1, sentence[offsets[i-3]:offsets[i]])
		}
		if i > 1 {
			score += lookupScore(p.tw2, sentence[offsets[i-2]:offsets[i+1]])
		}
		if i+1 < len(offsets)-1 {
			score += lookupScore(p.tw3, sentence[offsets[i-1]:offsets[i+2]])
		}
		if i+2 < len(offsets)-1 {
			score += lookupScore(p.tw4, sentence[offsets[i]:offsets[i+3]])
		}

		if score > 0 {
			boundaries = append(boundaries, offsets[i])
		}
	}
	boundaries = append(boundaries, len(sentence))

	result := make([]string, len(boundaries)-1)
	for i := 1; i < len(boundaries); i++ {
		start := boundaries[i-1]
		end := boundaries[i]
		result[i-1] = sentence[start:end]
	}
	return result
}

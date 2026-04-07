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
	scoreBias int

	// Valid UTF-8 feature keys are compiled into typed maps so Parse can avoid
	// allocating temporary strings and repeated string hashing on the hot path.
	uw1 map[rune]int
	uw2 map[rune]int
	uw3 map[rune]int
	uw4 map[rune]int
	uw5 map[rune]int
	uw6 map[rune]int

	// Invalid UTF-8 keys must stay byte-for-byte distinct to preserve the
	// historical behavior of New(models.Model) for custom models.
	uw1Raw map[string]int
	uw2Raw map[string]int
	uw3Raw map[string]int
	uw4Raw map[string]int
	uw5Raw map[string]int
	uw6Raw map[string]int

	bw1 map[uint64]int
	bw2 map[uint64]int
	bw3 map[uint64]int

	bw1Raw map[string]int
	bw2Raw map[string]int
	bw3Raw map[string]int

	tw1 map[uint64]int
	tw2 map[uint64]int
	tw3 map[uint64]int
	tw4 map[uint64]int

	tw1Raw map[string]int
	tw2Raw map[string]int
	tw3Raw map[string]int
	tw4Raw map[string]int
}

const runeKeyBits = 21

func packBigram(a, b rune) uint64 {
	return uint64(a)<<runeKeyBits | uint64(b)
}

func packTrigram(a, b, c rune) uint64 {
	return uint64(a)<<(runeKeyBits*2) | uint64(b)<<runeKeyBits | uint64(c)
}

type decodedKey struct {
	runes [3]rune
	count int
	valid bool
}

// decodeKey keeps enough information to choose between the fast typed-key path
// and the raw string fallback without losing invalid UTF-8 bytes.
func decodeKey(key string) decodedKey {
	decoded := decodedKey{valid: true}
	for len(key) > 0 {
		r, size := utf8.DecodeRuneInString(key)
		if r == utf8.RuneError && size == 1 {
			decoded.valid = false
		}
		if decoded.count < len(decoded.runes) {
			decoded.runes[decoded.count] = r
		}
		decoded.count++
		key = key[size:]
	}
	return decoded
}

func compileUnigramGroup(group map[string]int) (map[rune]int, map[string]int) {
	if len(group) == 0 {
		return nil, nil
	}

	var compiled map[rune]int
	var raw map[string]int
	for key, score := range group {
		decoded := decodeKey(key)
		if decoded.count != 1 {
			continue
		}
		if decoded.valid {
			if compiled == nil {
				compiled = make(map[rune]int, len(group))
			}
			compiled[decoded.runes[0]] = score
			continue
		}
		if raw == nil {
			raw = make(map[string]int, len(group))
		}
		raw[key] = score
	}
	return compiled, raw
}

func compileBigramGroup(group map[string]int) (map[uint64]int, map[string]int) {
	if len(group) == 0 {
		return nil, nil
	}

	var compiled map[uint64]int
	var raw map[string]int
	for key, score := range group {
		decoded := decodeKey(key)
		if decoded.count != 2 {
			continue
		}
		if decoded.valid {
			if compiled == nil {
				compiled = make(map[uint64]int, len(group))
			}
			compiled[packBigram(decoded.runes[0], decoded.runes[1])] = score
			continue
		}
		if raw == nil {
			raw = make(map[string]int, len(group))
		}
		raw[key] = score
	}
	return compiled, raw
}

func compileTrigramGroup(group map[string]int) (map[uint64]int, map[string]int) {
	if len(group) == 0 {
		return nil, nil
	}

	var compiled map[uint64]int
	var raw map[string]int
	for key, score := range group {
		decoded := decodeKey(key)
		if decoded.count != 3 {
			continue
		}
		if decoded.valid {
			if compiled == nil {
				compiled = make(map[uint64]int, len(group))
			}
			compiled[packTrigram(decoded.runes[0], decoded.runes[1], decoded.runes[2])] = score
			continue
		}
		if raw == nil {
			raw = make(map[string]int, len(group))
		}
		raw[key] = score
	}
	return compiled, raw
}

// Create new parser.
func New(model models.Model) *Parser {
	// BudouX uses -TotalScore/2 as the decision threshold.
	// Some bundled models have an odd total score, so rounding that threshold to
	// an integer would change segmentation results. Instead we scale the whole
	// comparison by 2 and compare against an integer bias: 2*score > 0.
	//
	// Feature lookups in Parse are on the hot path, so we compile the public
	// string-keyed model into typed keys here once and avoid repeated string
	// hashing for 1/2/3-rune feature windows. Keys that contain invalid UTF-8
	// stay in raw string fallback maps so New(models.Model) keeps its previous
	// byte-for-byte matching semantics for custom models.
	uw1, uw1Raw := compileUnigramGroup(model["UW1"])
	uw2, uw2Raw := compileUnigramGroup(model["UW2"])
	uw3, uw3Raw := compileUnigramGroup(model["UW3"])
	uw4, uw4Raw := compileUnigramGroup(model["UW4"])
	uw5, uw5Raw := compileUnigramGroup(model["UW5"])
	uw6, uw6Raw := compileUnigramGroup(model["UW6"])
	bw1, bw1Raw := compileBigramGroup(model["BW1"])
	bw2, bw2Raw := compileBigramGroup(model["BW2"])
	bw3, bw3Raw := compileBigramGroup(model["BW3"])
	tw1, tw1Raw := compileTrigramGroup(model["TW1"])
	tw2, tw2Raw := compileTrigramGroup(model["TW2"])
	tw3, tw3Raw := compileTrigramGroup(model["TW3"])
	tw4, tw4Raw := compileTrigramGroup(model["TW4"])

	return &Parser{
		scoreBias: -model.TotalScore(),
		uw1:       uw1,
		uw2:       uw2,
		uw3:       uw3,
		uw4:       uw4,
		uw5:       uw5,
		uw6:       uw6,
		uw1Raw:    uw1Raw,
		uw2Raw:    uw2Raw,
		uw3Raw:    uw3Raw,
		uw4Raw:    uw4Raw,
		uw5Raw:    uw5Raw,
		uw6Raw:    uw6Raw,
		bw1:       bw1,
		bw2:       bw2,
		bw3:       bw3,
		bw1Raw:    bw1Raw,
		bw2Raw:    bw2Raw,
		bw3Raw:    bw3Raw,
		tw1:       tw1,
		tw2:       tw2,
		tw3:       tw3,
		tw4:       tw4,
		tw1Raw:    tw1Raw,
		tw2Raw:    tw2Raw,
		tw3Raw:    tw3Raw,
		tw4Raw:    tw4Raw,
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

// Parses a sentence into phrases.
func (p *Parser) Parse(sentence string) []string {
	if sentence == "" {
		return []string{}
	}

	runeCount := utf8.RuneCountInString(sentence)
	if runeCount == 1 {
		return []string{sentence}
	}

	type runePos struct {
		r     rune
		start int
		end   int
		// valid reports whether this window element came from valid UTF-8 bytes.
		// Invalid bytes use the raw string fallback so custom model keys keep
		// their previous byte-oriented matching semantics.
		valid bool
		ok    bool
	}

	readNext := func(offset int) (runePos, int) {
		if offset >= len(sentence) {
			return runePos{}, offset
		}

		r, size := utf8.DecodeRuneInString(sentence[offset:])
		return runePos{
			r:     r,
			start: offset,
			end:   offset + size,
			valid: !(r == utf8.RuneError && size == 1),
			ok:    true,
		}, offset + size
	}

	offset := 0
	prev1, offset := readNext(offset)
	curr, offset := readNext(offset)
	next1, offset := readNext(offset)
	next2, offset := readNext(offset)
	prev2 := runePos{}
	prev3 := runePos{}

	boundaries := make([]int, 1, runeCount+1)

	for {
		score := p.scoreBias
		if prev3.ok {
			if prev3.valid {
				score += 2 * p.uw1[prev3.r]
			} else {
				// Keep invalid UTF-8 windows on the raw string path to avoid
				// collapsing distinct byte sequences into utf8.RuneError.
				score += 2 * p.uw1Raw[sentence[prev3.start:prev3.end]]
			}
		}
		if prev2.ok {
			if prev2.valid {
				score += 2 * p.uw2[prev2.r]
			} else {
				score += 2 * p.uw2Raw[sentence[prev2.start:prev2.end]]
			}
		}
		if prev1.valid {
			score += 2 * p.uw3[prev1.r]
		} else {
			score += 2 * p.uw3Raw[sentence[prev1.start:prev1.end]]
		}
		if curr.valid {
			score += 2 * p.uw4[curr.r]
		} else {
			score += 2 * p.uw4Raw[sentence[curr.start:curr.end]]
		}
		if next1.ok {
			if next1.valid {
				score += 2 * p.uw5[next1.r]
			} else {
				score += 2 * p.uw5Raw[sentence[next1.start:next1.end]]
			}
		}
		if next2.ok {
			if next2.valid {
				score += 2 * p.uw6[next2.r]
			} else {
				score += 2 * p.uw6Raw[sentence[next2.start:next2.end]]
			}
		}

		if prev2.ok {
			if prev2.valid && prev1.valid {
				score += 2 * p.bw1[packBigram(prev2.r, prev1.r)]
			} else {
				score += 2 * p.bw1Raw[sentence[prev2.start:prev1.end]]
			}
		}
		if prev1.valid && curr.valid {
			score += 2 * p.bw2[packBigram(prev1.r, curr.r)]
		} else {
			score += 2 * p.bw2Raw[sentence[prev1.start:curr.end]]
		}
		if next1.ok {
			if curr.valid && next1.valid {
				score += 2 * p.bw3[packBigram(curr.r, next1.r)]
			} else {
				score += 2 * p.bw3Raw[sentence[curr.start:next1.end]]
			}
		}

		if prev3.ok {
			if prev3.valid && prev2.valid && prev1.valid {
				score += 2 * p.tw1[packTrigram(prev3.r, prev2.r, prev1.r)]
			} else {
				score += 2 * p.tw1Raw[sentence[prev3.start:prev1.end]]
			}
		}
		if prev2.ok {
			if prev2.valid && prev1.valid && curr.valid {
				score += 2 * p.tw2[packTrigram(prev2.r, prev1.r, curr.r)]
			} else {
				score += 2 * p.tw2Raw[sentence[prev2.start:curr.end]]
			}
		}
		if next1.ok {
			if prev1.valid && curr.valid && next1.valid {
				score += 2 * p.tw3[packTrigram(prev1.r, curr.r, next1.r)]
			} else {
				score += 2 * p.tw3Raw[sentence[prev1.start:next1.end]]
			}
		}
		if next2.ok {
			if curr.valid && next1.valid && next2.valid {
				score += 2 * p.tw4[packTrigram(curr.r, next1.r, next2.r)]
			} else {
				score += 2 * p.tw4Raw[sentence[curr.start:next2.end]]
			}
		}

		if score > 0 {
			boundaries = append(boundaries, curr.start)
		}

		if !next1.ok {
			break
		}

		prev3 = prev2
		prev2 = prev1
		prev1 = curr
		curr = next1
		next1 = next2
		next2, offset = readNext(offset)
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

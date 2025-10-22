package budoux

import (
	"iter"

	"github.com/soundkitchen/go-budoux/models"
)

type Parser struct {
	model models.Model
}

// Create new parser.
func New(model models.Model) *Parser {
	return &Parser{
		model: model,
	}
}

// NewDefaultJapaneseParser returns new Parser with default japanese model.
func NewDefaultJapaneseParser() *Parser {
	return New(models.GetDefaultJapaneseModel())
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

// getScore returns the score of a word in a group if exists.
// otherwise returns defaultScore.
func (p *Parser) getScore(group string, word string, defaultScore float64) float64 {
	if g, ok := p.model[group]; ok {
		if score, ok := g[word]; ok {
			//return float64(score)
			return score
		}
	}
	return defaultScore
}

// Parses a sentence into phrases.
func (p *Parser) Parse(sentence string) []string {
	if sentence == "" {
		return []string{}
	}
	runes := []rune(sentence)
	singles := make([]string, len(runes))
	for i, r := range runes {
		singles[i] = string(r)
	}
	if len(runes) == 1 {
		return []string{singles[0]}
	}

	baseScore := -p.model.TotalScore() * 0.5
	boundaries := []int{0}

	for window := range slidingWindows(runes, singles) {
		score := baseScore
		if window.hasPrev3 {
			score += p.getScore("UW1", window.prev3, 0)
		}
		if window.hasPrev2 {
			score += p.getScore("UW2", window.prev2, 0)
		}
		score += p.getScore("UW3", window.prev1, 0)
		score += p.getScore("UW4", window.curr, 0)
		if window.hasNext1 {
			score += p.getScore("UW5", window.next1, 0)
		}
		if window.hasNext2 {
			score += p.getScore("UW6", window.next2, 0)
		}

		if window.hasPrev2Prev1 {
			score += p.getScore("BW1", window.prev2prev1, 0)
		}
		score += p.getScore("BW2", window.prev1curr, 0)
		if window.hasCurrNext {
			score += p.getScore("BW3", window.currnext, 0)
		}

		if window.hasPrev3Prev2Prev1 {
			score += p.getScore("TW1", window.prev3prev2prev1, 0)
		}
		if window.hasPrev2Prev1Curr {
			score += p.getScore("TW2", window.prev2prev1curr, 0)
		}
		if window.hasPrev1CurrNext1 {
			score += p.getScore("TW3", window.prev1currnext1, 0)
		}
		if window.hasCurrNext1Next2 {
			score += p.getScore("TW4", window.currnext1next2, 0)
		}

		if score > 0 {
			boundaries = append(boundaries, window.index)
		}
	}
	boundaries = append(boundaries, len(runes))

	result := make([]string, len(boundaries)-1)
	for i := 1; i < len(boundaries); i++ {
		start := boundaries[i-1]
		end := boundaries[i]
		result[i-1] = string(runes[start:end])
	}
	return result
}

type runeWindow struct {
	index int

	prev3    string
	hasPrev3 bool

	prev2    string
	hasPrev2 bool

	prev1 string

	curr string

	next1    string
	hasNext1 bool

	next2    string
	hasNext2 bool

	prev2prev1    string
	hasPrev2Prev1 bool

	prev1curr string

	currnext    string
	hasCurrNext bool

	prev3prev2prev1    string
	hasPrev3Prev2Prev1 bool

	prev2prev1curr    string
	hasPrev2Prev1Curr bool

	prev1currnext1    string
	hasPrev1CurrNext1 bool

	currnext1next2    string
	hasCurrNext1Next2 bool
}

func slidingWindows(runes []rune, singles []string) iter.Seq[runeWindow] {
	return func(yield func(runeWindow) bool) {
		if len(runes) < 2 {
			return
		}

		for i := 1; i < len(runes); i++ {
			window := runeWindow{
				index:     i,
				prev1:     singles[i-1],
				curr:      singles[i],
				prev1curr: string(runes[i-1 : i+1]),
			}

			if i > 2 {
				window.prev3 = singles[i-3]
				window.hasPrev3 = true
				window.prev3prev2prev1 = string(runes[i-3 : i])
				window.hasPrev3Prev2Prev1 = true
			}

			if i > 1 {
				window.prev2 = singles[i-2]
				window.hasPrev2 = true
				window.prev2prev1 = string(runes[i-2 : i])
				window.hasPrev2Prev1 = true
				window.prev2prev1curr = string(runes[i-2 : i+1])
				window.hasPrev2Prev1Curr = true
			}

			if i+1 < len(runes) {
				window.next1 = singles[i+1]
				window.hasNext1 = true
				window.currnext = string(runes[i : i+2])
				window.hasCurrNext = true
				window.prev1currnext1 = string(runes[i-1 : i+2])
				window.hasPrev1CurrNext1 = true
			}

			if i+2 < len(runes) {
				window.next2 = singles[i+2]
				window.hasNext2 = true
				window.currnext1next2 = string(runes[i : i+3])
				window.hasCurrNext1Next2 = true
			}

			if !yield(window) {
				return
			}
		}
	}
}

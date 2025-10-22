package models

import "maps"

type Model map[string]map[string]float64

// TotalScore calculates the total score of the model.
func (m Model) TotalScore() float64 {
	var total float64
	for group := range maps.Values(m) {
		for score := range maps.Values(group) {
			total += score
		}
	}
	return total
}

func GetDefaultJapaneseModel() Model {
	return ja
}

func GetDefaultThaiModel() Model {
	return th
}

func GetDefaultSimplifiedChineseModel() Model {
	return zhhans
}

func GetDefaultTraditionalChineseModel() Model {
	return zhhant
}

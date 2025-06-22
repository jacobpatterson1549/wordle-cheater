package server

import (
	"slices"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/spelling_bee"
)

type (
	SpellingBeeCheater struct {
		spelling_bee.SpellingBee
		Output []string
	}
	Word struct {
		Score     int
		Value     string
		IsPangram string
	}
)

func RunSpellingBeeCheater(query map[string][]string) SpellingBeeCheater {
	var sbc SpellingBeeCheater
	sbc.MinLength = 4

	centralLetters, ok := query["central-letter"]
	switch {
	case !ok:
		// NOOP
	case len(centralLetters) != 1:
		sbc.Output = append(sbc.Output, "only one central letter parameter allowed")
	case len(centralLetters[0]) != 1:
		sbc.Output = append(sbc.Output, "central letter must be one character")
	default:
		sbc.CentralLetter = []rune(centralLetters[0])[0]
	}

	otherLetters, ok := query["other-letters"]
	switch {
	case !ok:
		// NOOP
	case len(otherLetters) != 1:
		sbc.Output = append(sbc.Output, "only one other letters parameter allowed")
	case len(otherLetters[0]) != 6:
		sbc.Output = append(sbc.Output, "other letters must be six characters long")
	default:
		sbc.OtherLetters = otherLetters[0]
	}

	return sbc
}

func (sbc SpellingBeeCheater) Words() []Word {
	words := sbc.SpellingBee.Words(words.WordsTextFile)
	display := make([]Word, len(words))
	for i, w := range words {
		display[i] = Word{
			Score: w.Score,
			Value: w.Value,
		}
		if w.IsPangram {
			display[i].IsPangram = "PANGRAM!"
		}
	}
	slices.Reverse(display)
	return display
}

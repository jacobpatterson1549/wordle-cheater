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
	Summary struct {
		TotalScore   int
		PangramCount int
		Words        []Word
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

func (sbc SpellingBeeCheater) Summary() Summary {
	var s Summary
	words := sbc.SpellingBee.Words(words.WordsTextFile)
	s.Words = make([]Word, len(words))
	for i, w := range words {
		s.Words[i].Value = w.Value
		s.Words[i].Score = w.Score
		s.TotalScore += w.Score
		if w.IsPangram {
			s.Words[i].IsPangram = "PANGRAM!"
			s.PangramCount++
		}
	}
	slices.Reverse(s.Words)
	return s
}

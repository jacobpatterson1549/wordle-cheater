package server

import (
	"cmp"
	"fmt"
	"slices"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/spelling_bee"
)

type (
	SpellingBeeCheater struct {
		spelling_bee.SpellingBee
	}
	Summary struct {
		TotalScore   int
		PangramCount int
		Words        []Word
	}
	Word struct {
		Score   int
		Value   string
		Details string
	}
)

const (
	centralLetterParam = "central-letter"
	otherLettersParam  = "other-letters"
)

func RunSpellingBeeCheater(query map[string][]string) (*SpellingBeeCheater, error) {
	centralLetters, err1 := parseParam(centralLetterParam, 1, query)
	otherLetters, err2 := parseParam(otherLettersParam, 6, query)
	if err := cmp.Or(err1, err2); err != nil {
		return nil, err
	}
	if (len(centralLetters) == 0) != (len(otherLetters) == 0) {
		return nil, fmt.Errorf("both params not specified")
	}
	sb := spelling_bee.SpellingBee{
		OtherLetters: otherLetters,
		MinLength:    4,
	}
	for _, r := range centralLetters {
		sb.CentralLetter = r
	}
	return &SpellingBeeCheater{sb}, nil
}

func parseParam(paramName string, wantLength int, query map[string][]string) (string, error) {
	value, ok := query[paramName]
	switch {
	case !ok:
		return "", nil
	case len(value) != 1:
		return "", fmt.Errorf("only one %q parameter allowed", paramName)
	case len(value[0]) != wantLength:
		return "", fmt.Errorf("%q must be %v characters long", paramName, wantLength)
	}
	return value[0], nil
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
			s.Words[i].Details = "PANGRAM!"
			s.PangramCount++
		}
	}
	slices.Reverse(s.Words)
	return s
}

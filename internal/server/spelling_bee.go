package server

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/jacobpatterson1549/wordle-cheater/internal/spelling_bee"
)

type (
	SpellingBeeCheater struct {
		spelling_bee.SpellingBee
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

func NewSpellingBeeCheater(query map[string][]string, wordsText string) (*SpellingBeeCheater, error) {
	sb, err := newSpellingBee(query)
	if err != nil {
		return nil, err
	}
	sbc := newSpellingBeeCheater(*sb, wordsText)
	return sbc, nil
}

func newSpellingBee(query map[string][]string) (*spelling_bee.SpellingBee, error) {
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
	return &sb, nil
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

func newSpellingBeeCheater(sb spelling_bee.SpellingBee, wordsText string) *SpellingBeeCheater {
	sbc := SpellingBeeCheater{
		SpellingBee: sb,
	}
	words := sb.Words(wordsText)
	sbc.Words = make([]Word, len(words))
	for i, w := range words {
		sbc.Words[i].Value = w.Value
		sbc.Words[i].Score = w.Score
		sbc.TotalScore += w.Score
		if w.IsPangram {
			sbc.Words[i].Details = "PANGRAM!"
			sbc.PangramCount++
		}
	}
	slices.Reverse(sbc.Words)
	return &sbc
}

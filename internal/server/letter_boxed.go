package server

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/jacobpatterson1549/wordle-cheater/internal/letter_boxed"
)

type (
	LetterBoxedCheater struct {
		letter_boxed.LetterBox
		letter_boxed.Result
	}
)

const (
	letterBoxedLettersParam = "letters"
)

func NewLetterBoxedCheater(query map[string][]string, wordsText string) (*LetterBoxedCheater, error) {
	lb, err := newLetterBox(query)
	if err != nil {
		return nil, fmt.Errorf("parsing params: %w", err)
	}
	r, err := lb.Solve(wordsText)
	if err != nil {
		return nil, fmt.Errorf("searching for words: %v", err)
	}
	lbc := LetterBoxedCheater{
		LetterBox: *lb,
		Result:    *r,
	}
	slices.SortFunc(lbc.Result.Words, lbc.sortWords)
	return &lbc, nil
}

func newLetterBox(query map[string][]string) (*letter_boxed.LetterBox, error) {
	lb := letter_boxed.LetterBox{
		BoxSideCount:  4,
		MinWordLength: 3,
	}
	letters := query[letterBoxedLettersParam]
	switch n := len(letters); {
	case n > 1:
		return nil, fmt.Errorf("wanted only one %q parameter, got %v", letterBoxedLettersParam, n)
	case n == 1:
		lb.Letters = letters[0]
	}
	return &lb, nil
}

func (LetterBoxedCheater) sortWords(a, b string) int {
	if len(a) != len(b) {
		return len(b) - len(a)
	}
	return cmp.Compare(a, b)
}

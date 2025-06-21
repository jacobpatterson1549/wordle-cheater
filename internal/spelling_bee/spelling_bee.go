package spelling_bee

import (
	"slices"
	"strings"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/char_set"
)

type (
	SpellingBee struct {
		CentralLetter rune
		OtherLetters  string
		MinLength     int
		validLetters  char_set.CharSet
		numLetters    int
	}
	Word struct {
		Value     string
		Score     int
		IsPangram bool
	}
)

func (sb SpellingBee) Words(wordsTextFile string) []Word {
	lines := strings.Fields(wordsTextFile)
	sb.init()
	var words []Word
	for _, value := range lines {
		letters := sb.letters(value)
		if letters != 0 {
			w := sb.newWord(value, letters)
			words = append(words, w)
		}
	}
	slices.SortFunc(words, wordLess)
	return words
}

func (sb *SpellingBee) init() {
	sb.validLetters = 0
	if !lowercase(sb.CentralLetter) {
		return
	}
	sb.validLetters.Add(sb.CentralLetter)
	for _, r := range sb.OtherLetters {
		if lowercase(r) {
			sb.validLetters.Add(r)
		}
	}
	sb.numLetters = sb.validLetters.Length()
}

func (sb SpellingBee) letters(value string) char_set.CharSet {
	if len(value) < sb.MinLength {
		return 0
	}
	var letters char_set.CharSet
	for _, r := range value {
		if !sb.validLetters.Has(r) {
			return 0
		}
		letters.Add(r)
	}
	if !letters.Has(sb.CentralLetter) {
		return 0
	}
	return letters
}

func (sb *SpellingBee) newWord(value string, letters char_set.CharSet) Word {
	w := Word{
		Value:     value,
		Score:     1,
		IsPangram: letters == sb.validLetters,
	}
	if sb.MinLength < len(value) {
		w.Score += len(value) - 1
	}
	if w.IsPangram {
		w.Score += sb.numLetters
	}
	return w
}

func lowercase(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func wordLess(a, b Word) int {
	switch {
	case a.Score != b.Score:
		return a.Score - b.Score
	case a.IsPangram != b.IsPangram && b.IsPangram:
		return -1
	case a.IsPangram != b.IsPangram && a.IsPangram:
		return 1
	case len(a.Value) != len(b.Value):
		return len(a.Value) - len(b.Value)
	case a.Value != b.Value:
		return strings.Compare(a.Value, b.Value)
	}
	return 0
}

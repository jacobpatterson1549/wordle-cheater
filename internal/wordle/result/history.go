package result

import (
	"fmt"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/char_set"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/guess"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/score"
)

type (
	// History stores the state of multiple results
	History struct {
		correctLetters    [numLetters]rune
		almostLetters     []rune
		prohibitedLetters [numLetters]char_set.CharSet
	}
	// Result is a guess and it's score
	Result struct {
		Guess guess.Guess
		Score score.Score
	}
)

const numLetters = 5

// addResult merges the result into the history and trims the words to only include ones that are allowed
func (h *History) AddResult(r Result, m *words.Words) {
	h.mergeResult(r)
	for w := range *m {
		if !h.allows(w) {
			delete(*m, w)
		}
	}
}

// mergeResult merges the result into the history
func (h *History) mergeResult(r Result) {
	var usedLetters []rune
	for i, si := range r.Score {
		gi := rune(r.Guess[i])
		switch si {
		case 'c':
			h.setLetterCorrect(gi, i)
			usedLetters = append(usedLetters, gi)
		case 'a':
			h.setLetterAlmost(gi, i)
			usedLetters = append(usedLetters, gi)
		case 'n':
			h.setLetterProhibited(gi, i, usedLetters)
		}
	}
	h.mergeRequiredLetters(usedLetters)
}

// setLetterCorrect sets the letter at the index to correct
func (h *History) setLetterCorrect(ch rune, index int) {
	h.correctLetters[index] = ch
}

// setLetterAlmost marks the letter as available somewhere else by prohibiting it at the index
func (h *History) setLetterAlmost(ch rune, index int) {
	p := &h.prohibitedLetters[index]
	p.Add(ch)
}

// setLetterProhibited marks the letter as prohibited from all indexes
func (h *History) setLetterProhibited(ch rune, index int, usedLetters []rune) {
	for j := 0; j < numLetters; j++ {
		pj := &h.prohibitedLetters[j]
		pj.Add(ch)
	}
}

// mergeRequiredLetters adds required letters from a guess into the required letters.
// New letters are only added if they were not previously required.
func (h *History) mergeRequiredLetters(usedLetters []rune) {
	requiredLetters := make([]rune, len(h.almostLetters))
	copy(requiredLetters, h.almostLetters)
	existingCounts := letterCounts(h.almostLetters...)
	scoreCounts := letterCounts(usedLetters...)
	for _, ch := range usedLetters {
		if existingCounts[ch] < scoreCounts[ch] {
			scoreCounts[ch]--
			requiredLetters = append(requiredLetters, ch)
		}
	}
	h.almostLetters = requiredLetters
}

// letterCounts creates a count multi-map of the runes (count-set)
func letterCounts(runes ...rune) map[rune]int {
	m := make(map[rune]int, len(runes))
	for _, r := range runes {
		m[r]++
	}
	return m
}

// allows determines if a word is allowed based on the history (not prohibited)
func (h *History) allows(w string) bool {
	letterCounts := make(map[rune]int, numLetters)
	for i, ch := range w {
		switch {
		case h.correctLetters[i] != 0 && h.correctLetters[i] != ch,
			h.correctLetters[i] == 0 && h.prohibitedLetters[i].Has(ch):
			return false
		}
		letterCounts[ch]++
	}
	for _, ch := range h.almostLetters {
		n, ok := letterCounts[ch]
		switch {
		case !ok:
			return false // required letter not present
		case n == 1:
			delete(letterCounts, ch)
		default:
			letterCounts[ch]--
		}
	}
	return true
}

// String formats the required and prohibited letters to clearly show the state
func (h History) String() string {
	correct := make([]rune, len(h.correctLetters))
	for i, ch := range h.correctLetters {
		switch {
		case ch == 0:
			correct[i] = '?'
		default:
			correct[i] = ch
		}
	}
	almost := make([]string, len(h.almostLetters))
	for i, ch := range h.almostLetters {
		almost[i] = string(ch)
	}
	prohibited := make([]string, len(h.prohibitedLetters))
	for i, cs := range h.prohibitedLetters {
		prohibited[i] = cs.String()
	}
	a := struct {
		correctLetters    string
		almostLetters     []string
		prohibitedLetters []string
	}{
		string(correct),
		almost,
		prohibited,
	}
	return fmt.Sprintf("%+v", a)
}

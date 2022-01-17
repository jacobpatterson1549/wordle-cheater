// Package main runs a command-line-interface program to help cheat in the popular Worlde game
package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

//go:embed build/words.txt
var wordsTextFile string

// numLetters is the length of the words
const numLetters = 5

// main runs wordle-cheater on the command-line using stdin and stdout
func main() {
	var rw osReadWriter
	if err := runWordleCheater(rw, wordsTextFile); err != nil {
		panic(fmt.Errorf("running wordle: %v", err))
	}
}

// runWordleCheater runs an interactive wordle-cheater on the ReaderWriter using the text for the words
func runWordleCheater(rw io.ReadWriter, wordsText string) error {
	allWords, err := newWords(wordsText)
	if err != nil {
		return fmt.Errorf("loading words: %v", err)
	}
	availableWords := allWords.copy()

	fmt.Fprintf(rw, "Running wordle-cheater\n")
	fmt.Fprintf(rw, " * Guesses and scores are %v letters long\n", numLetters)
	fmt.Fprintf(rw, " * Scores are only made of the following letters:\n")
	fmt.Fprintf(rw, "   C - if a letter is in the word and in the correct location\n")
	fmt.Fprintf(rw, "   A - if a letter is in the word, but in the wrong location\n")
	fmt.Fprintf(rw, "   N - if a letter is not in the word\n")
	fmt.Fprintf(rw, "The app runs until the correct word is found from a guess with only correct letters.\n\n")

	var h history
	for {
		g, err := newGuess(rw, *allWords)
		if err != nil {
			return err
		}

		s, err := newScore(rw)
		if err != nil {
			return err
		}
		if s.allCorrect() {
			return nil
		}

		r := result{
			guess: *g,
			score: *s,
		}
		if err := h.addResult(r, availableWords); err != nil {
			return err
		}

		if err := availableWords.scanShowPossible(rw); err != nil {
			return err
		}
	}
}

// osReadWriter wraps reading from the standard input and output, satisfying io.ReadWriter
type osReadWriter struct{}

// Read reads p from stdin
func (rw osReadWriter) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

// Write writes p to stdout
func (rw osReadWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// words is a collection of unique strings
type words map[string]struct{}

// newWords loads the words from the file
// words are separated by whitespace (spaces/newlines)
// an error is returned if any words are not <<numLetters characters long and lowercase
func newWords(a string) (*words, error) {
	lines := strings.Fields(a)
	m := make(words, len(lines))
	for _, w := range lines {
		if len(w) != numLetters {
			return nil, fmt.Errorf("wanted all words to be %v letters long, got %q", numLetters, w)
		}
		if w != strings.ToLower(w) {
			return nil, fmt.Errorf("wanted all words to be lowercase, got %q", w)
		}
		m[w] = struct{}{}
	}
	return &m, nil
}

// copy creates a new, identical duplication of the words
func (m words) copy() *words {
	m2 := make(words, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return &m2
}

// sorted combines and sorts the words into a csv string
func (m words) sorted() string {
	s := make([]string, 0, len(m))
	for w := range m {
		s = append(s, w)
	}
	sort.Strings(s)
	j := strings.Join(s, ",")
	return j
}

// scanShowPossible prompts to display the words
func (m words) scanShowPossible(rw io.ReadWriter) error {
	fmt.Fprintf(rw, "show possible words [Yn]: ")
	var choice string
	if _, err := fmt.Fscan(rw, &choice); err != nil {
		return fmt.Errorf("scanning choice: %v", err)
	}
	choice = strings.ToLower(choice)
	if len(choice) > 0 && choice[0] != 'y' {
		return nil
	}
	fmt.Fprintf(rw, "remaining valid words: %v\n", m.sorted())
	return nil
}

// guess is a word that might be the answer
type guess string

// newGuess prompts for a guess on the ReadWriter until a valid one is given or an io error occurs
func newGuess(rw io.ReadWriter, m words) (*guess, error) {
	for {
		fmt.Fprintf(rw, "Enter guess (%v letters): ", numLetters)
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		word = strings.ToLower(word)
		g := guess(word)
		if err := g.validate(m); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &g, nil
	}
}

// validate ensures the guess is <<numLetters>> letters long and is in the words list (if a list is provided)
func (g guess) validate(m words) error {
	if len(g) != numLetters {
		return fmt.Errorf("guess must be %v letters long", numLetters)
	}
	if len(m) > 0 {
		if _, ok := m[string(g)]; !ok {
			return fmt.Errorf("%v is not a word", g)
		}
	}
	return nil
}

// score is a <<numLetters>>-letter string made up of {c,a,n}
// The letter c indicates that a letter from a guess is in the correct position
// The letter a indicates that a letter from a guess is in the answer, but in a different position
// The letter n indicates that a letter from a guess is not anywhere in the answer
type score string

// newScore prompts for a score on the ReadWriter until a valid one is given or an io error occurs
func newScore(rw io.ReadWriter) (*score, error) {
	for {
		fmt.Fprintf(rw, "Enter score: ")
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		word = strings.ToLower(word)
		s := score(word)
		if err := s.validate(); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &s, nil
	}
}

// validate ensures the score is <<numLeters>> letters long and consists only of the {c,a,n} letters
func (s score) validate() error {
	if len(s) != numLetters {
		return fmt.Errorf("score must be %v letters long", numLetters)
	}
	for _, ch := range s {
		switch ch {
		case 'c', 'a', 'n':
			// NOOP
		default:
			return fmt.Errorf("must be only the following letters: C, A, N")
		}
	}
	return nil
}

// allCorrect determines if every character in the score is c
func (s score) allCorrect() bool {
	return s == "ccccc"
}

// result is a guess and it's score
type result struct {
	guess guess
	score score
}

// validate returns an error if the guess or score is invalid
func (r result) validate() error {
	if err := r.guess.validate(nil); err != nil {
		return fmt.Errorf("validating guess: %v", err)
	}
	if err := r.score.validate(); err != nil {
		return fmt.Errorf("validating score: %v", err)
	}
	return nil
}

// history stores the state of multiple results
type history struct {
	requiredLetters   []rune
	prohibitedLetters [numLetters]charSet
}

// addResult merges the result into the history and trims the words to only include ones that are allowed
func (h *history) addResult(r result, m *words) error {
	if err := r.validate(); err != nil {
		return fmt.Errorf("adding invalid result to history: %v", err)
	}
	if err := h.mergeResult(r); err != nil {
		return fmt.Errorf("merging score: %v", err)
	}
	for w := range *m {
		if !h.allows(w) {
			delete(*m, w)
		}
	}
	return nil
}

// mergeResult merges the result into the history
func (h *history) mergeResult(r result) error {
	var usedLetters []rune
	for i, si := range r.score {
		gi := rune(r.guess[i])
		switch si {
		case 'c':
			if err := h.setLetterCorrect(gi, i); err != nil {
				return err
			}
			usedLetters = append(usedLetters, gi)
		case 'a':
			if err := h.setLetterAlmost(gi, i); err != nil {
				return err
			}
			usedLetters = append(usedLetters, gi)
		case 'n':
			if err := h.setLetterProhibited(gi, i, usedLetters); err != nil {
				return err
			}
		}
	}
	if err := h.mergeRequiredLetters(usedLetters...); err != nil {
		return fmt.Errorf("merging required letters: %v", err)
	}
	return nil
}

// setLetterCorrect sets the letter at the index to correct
func (h *history) setLetterCorrect(ch rune, index int) error {
	p := &h.prohibitedLetters[index]
	if p.Has(ch) {
		return fmt.Errorf("%c was prohibited at index %v, but is now supposedly correct", ch, index)
	}
	for l := 'a'; l <= 'z'; l++ {
		if l != ch {
			p.Add(l)
		}
	}
	return nil
}

// setLetterAlmost markes the letter as available somewhere else by prohibiting it at the index
func (h *history) setLetterAlmost(ch rune, index int) error {
	p := &h.prohibitedLetters[index]
	if p.AddWouldFill(ch) {
		return fmt.Errorf("all letters prohibited at index %v", index)
	}
	p.Add(ch)
	return nil
}

// setLetterProhibited marks the letter as prohibited from all indexes
func (h *history) setLetterProhibited(ch rune, index int, usedLetters []rune) error {
	if h.hasRequiredLetter(ch, usedLetters...) {
		return fmt.Errorf("%c was previously required to be in word, but is prohibited", ch)
	}
	for j := 0; j < numLetters; j++ {
		pj := &h.prohibitedLetters[j]
		if pj.AddWouldFill(ch) {
			return fmt.Errorf("all letters prohibited at index %v", index)
		}
		pj.Add(ch)
	}
	return nil
}

// hasRequiredLetter determines if the first rune is in requiredLetters or any of newRequiredLetters
func (h history) hasRequiredLetter(r rune, newRequiredLetters ...rune) bool {
	for _, ch := range h.requiredLetters {
		if r == ch {
			return true
		}
	}
	for _, ch := range newRequiredLetters {
		if r == ch {
			return true
		}
	}
	return false
}

// mergeRequiredLetters adds required letters from a guess into the required letters
// new letters are only added if they were not previously required
// an error is returned if too many letters are required
func (h *history) mergeRequiredLetters(newScoreLetters ...rune) error {
	requiredLetters := make([]rune, len(h.requiredLetters))
	copy(requiredLetters, h.requiredLetters)
	existingCounts := letterCounts(h.requiredLetters...)
	scoreCounts := letterCounts(newScoreLetters...)
	for _, ch := range newScoreLetters {
		if existingCounts[ch] < scoreCounts[ch] {
			scoreCounts[ch]--
			requiredLetters = append(requiredLetters, ch)
		}
	}
	if len(requiredLetters) > numLetters {
		return fmt.Errorf("more than five letters are now required")
	}
	h.requiredLetters = requiredLetters
	return nil
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
func (h *history) allows(w string) bool {
	letterCounts := make(map[rune]int, numLetters)
	for i, ch := range w {
		if h.prohibitedLetters[i].Has(ch) {
			return false
		}
		letterCounts[ch]++
	}
	for _, ch := range h.requiredLetters {
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

// charSet is a bitflag that stores the letters a-z
type charSet uint32

// Add includes the character to the set, panicing if the character is not in a-z
func (cs *charSet) Add(ch rune) {
	if !cs.valid(ch) {
		panic(fmt.Errorf("%c is not in a-z", ch))
	}
	*cs |= cs.singleton(ch)
}

// Has determines if the character is in the set
func (cs charSet) Has(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return (cs & cs.singleton(ch)) != 0
}

// AddWouldFill determines if the charset is filled with the letters a-z
func (cs charSet) AddWouldFill(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return cs|cs.singleton(ch) == cs.singleton('z'+1)-1
}

// String creates a string of the characters in the set, in ascending order
func (cs charSet) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for ch := rune('a'); ch <= 'z'; ch++ {
		if cs.Has(ch) {
			b.WriteRune(ch)
		}
	}
	b.WriteRune(']')
	return b.String()
}

// valid determines if the byte can be used in the charSet, if it is a-z
func (charSet) valid(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}

// singleton creates a singleton charSet from the character
func (charSet) singleton(ch rune) charSet {
	return 1 << (ch - 'a')
}

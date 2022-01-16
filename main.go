package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

//go:embed build/words.txt
var wordsTextFile string

func main() {
	var rw osReadWriter
	if err := runWordle(rw, wordsTextFile); err != nil {
		log.Fatalf("running wordle: %v", err)
	}
}

type osReadWriter struct{}

func (rw osReadWriter) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (rw osReadWriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func runWordle(rw io.ReadWriter, wordsLines string) error {
	words, err := newWords(wordsLines)
	if err != nil {
		return fmt.Errorf("loading words: %v", err)
	}
	fmt.Fprintf(rw, "There are %v legal words\n", len(*words)) // There are 2315 legal words

	fmt.Fprintln(rw, "Running wordle-cheater")
	fmt.Fprintln(rw, " * Guesses and scores are 5 letters long")
	fmt.Fprintln(rw, " * Scores are only made of the following letters:")
	fmt.Fprintln(rw, "   C - if a letter is in the word and in the correct location")
	fmt.Fprintln(rw, "   A - if a letter is in the word, but in the wrong location")
	fmt.Fprintln(rw, "   N - if a letter is not in the word")
	fmt.Fprintln(rw, "The app runs until the correct word is found from a guess with only correct letters.")
	fmt.Fprintln(rw)

	var h history
	for {
		g, err := newGuess(rw, *words)
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
		if err := h.addResult(r, words); err != nil {
			return err
		}

		if err := words.scanShowPossible(rw); err != nil {
			return err
		}
	}
}

type words map[string]struct{}

func newWords(a string) (*words, error) {
	lines := strings.Split(a, "\n")
	words := make(words, len(lines))
	for _, w := range lines {
		if len(w) != 5 {
			return nil, fmt.Errorf("wanted all words to be 5 letters long, got %q", w)
		}
		if w != strings.ToLower(w) {
			return nil, fmt.Errorf("wanted all words to be lowercase, got %q", w)
		}
		words[w] = struct{}{}
	}
	return &words, nil
}

func (words words) sorted() string {
	s := make([]string, len(words))
	for w := range words {
		s = append(s, w)
	}
	sort.Strings(s)
	j := strings.Join(s, ",")
	return j
}

func (words words) scanShowPossible(rw io.ReadWriter) error {
	fmt.Fprintf(rw, "show possible words [Yn]: ")
	var choice string
	if _, err := fmt.Fscan(rw, &choice); err != nil {
		return fmt.Errorf("scanning choice: %v", err)
	}
	choice = strings.ToLower(choice)
	if len(choice) > 0 && choice[0] != 'y' {
		return nil
	}
	fmt.Fprintf(rw, "remaining valid words: %v\n", words.sorted())
	return nil
}

type guess string

func newGuess(rw io.ReadWriter, words words) (*guess, error) {
	for {
		fmt.Fprintf(rw, "Enter guess (five letters): ")
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		word = strings.ToLower(word)
		g := guess(word)
		if err := g.validate(words); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &g, nil
	}
}

func (g guess) validate(words words) error {
	if n := 5; len(g) != n {
		return fmt.Errorf("guess must be %v letters long", n)
	}
	if len(words) > 0 {
		if _, ok := words[string(g)]; !ok {
			return fmt.Errorf("%v is not a word", g)
		}
	}
	return nil
}

type score string

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

func (s score) validate() error {
	if n := 5; len(s) != n {
		return fmt.Errorf("score must be %v letters long", n)
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

func (s score) allCorrect() bool {
	return s == "ccccc"
}

type result struct {
	guess guess
	score score
}

func (r result) validate() error {
	if err := r.guess.validate(nil); err != nil {
		return fmt.Errorf("validating guess: %v", err)
	}
	if err := r.score.validate(); err != nil {
		return fmt.Errorf("validating score: %v", err)
	}
	return nil
}

type history struct {
	results            []result
	requiredLetters    []rune
	prohibitedLetters1 charSet
	prohibitedLetters2 charSet
	prohibitedLetters3 charSet
	prohibitedLetters4 charSet
	prohibitedLetters5 charSet
}

func (h *history) addResult(r result, words *words) error {
	if err := r.validate(); err != nil {
		return fmt.Errorf("adding invalid result to history: %v", err)
	}
	if err := h.mergeResult(r); err != nil {
		return fmt.Errorf("merging score: %v", err)
	}
	h.results = append(h.results, r)
	for w := range *words {
		if !h.allows(w) {
			delete(*words, w)
		}
	}
	return nil
}

func (h *history) mergeResult(r result) error {
	var usedLetters []rune
	for i, si := range r.score {
		gi := rune(r.guess[i])
		p := h.prohibitedLetters(i)
		switch si {
		case 'c':
			if p.has(gi) {
				return fmt.Errorf("%c was prohibited at inde %v, but is now supposedly correct", si, i)
			}
			for l := 'a'; l <= 'z'; l++ {
				if l != gi {
					p.add(l)
				}
			}
			usedLetters = append(usedLetters, gi)
		case 'a':
			p.add(gi)
			if p.isFull() {
				return fmt.Errorf("all letters prohibited at index %v", i)
			}
			usedLetters = append(usedLetters, gi)
		case 'n':
			if h.hasRequiredLetter(gi, usedLetters...) {
				return fmt.Errorf("%c was previously required to be in word, but is prohibited", gi)
			}
			for j := range r.score {
				pj := h.prohibitedLetters(j)
				pj.add(gi)
				if pj.isFull() {
					return fmt.Errorf("all letters prohibited at index %v", i)
				}
			}
		}
	}
	if err := h.mergeRequiredLetters(usedLetters...); err != nil {
		return fmt.Errorf("merging required letters: %v", err)
	}
	return nil
}

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

func (h *history) mergeRequiredLetters(newScoreLetters ...rune) error {
	existingCounts := letterCounts(h.requiredLetters...)
	scoreCounts := letterCounts(newScoreLetters...)
	for _, ch := range newScoreLetters {
		if existingCounts[ch] < scoreCounts[ch] {
			scoreCounts[ch]--
			h.requiredLetters = append(h.requiredLetters, ch)
		}
	}
	if len(h.requiredLetters) > 5 {
		return fmt.Errorf("more than five letters are now required")
	}
	return nil
}

func letterCounts(runes ...rune) map[rune]int {
	m := make(map[rune]int, len(runes))
	for _, r := range runes {
		m[r]++
	}
	return m
}

func (h *history) prohibitedLetters(index int) *charSet {
	switch index {
	case 0:
		return &h.prohibitedLetters1
	case 1:
		return &h.prohibitedLetters2
	case 2:
		return &h.prohibitedLetters3
	case 3:
		return &h.prohibitedLetters4
	case 4:
		return &h.prohibitedLetters5
	default:
		panic(fmt.Errorf("unknown prohibited letter index: %v", index))
	}
}

func (h *history) allows(w string) bool {
	letterCounts := make(map[rune]int, 5)
	for i, ch := range w {
		if h.prohibitedLetters(i).has(ch) {
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

// add includes the character to the set, panicing if the character is not in a-z
func (cs *charSet) add(ch rune) {
	if !cs.valid(ch) {
		panic(fmt.Errorf("%c is not in a-z", ch))
	}
	*cs |= cs.singleton(ch)
}

// has determines if the character is in the set
func (cs charSet) has(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return (cs & cs.singleton(ch)) != 0
}

// isFull determines if the charset is filled with the letters a-z
func (cs charSet) isFull() bool {
	return cs == cs.singleton('z'+1)-1
}

// String creates a string of the characters in the set, in ascending order
func (cs charSet) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for ch := rune('a'); ch <= 'z'; ch++ {
		if cs.has(ch) {
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

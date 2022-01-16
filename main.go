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

//go:embed words.txt
var wordsTextFile string

// nasty  nannc
// ready  naanc
// early  aannc
// great  nnaan
// abbey  ccccc
/*
Enter guess (five letters): nasty
Enter score: nannc
show possible words [Yn]: y
remaining valid words: [basic lilac magic sumac iliac havoc antic attic]
Enter guess (five letters): ready
Enter score: naanc
show possible words [Yn]: y
remaining valid words: [basic lilac magic sumac iliac havoc antic attic]
Enter guess (five letters): early
Enter score: aannc
show possible words [Yn]: y
remaining valid words: [basic lilac magic sumac iliac havoc antic attic]
Enter guess (five letters): ^Csignal: interrupt

*/

func main() {
	var rw osReadWriter
	if err := runWordle(rw); err != nil {
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

func runWordle(rw io.ReadWriter) error {
	words, err := newWords(wordsTextFile)
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

	h := newHistory()
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
	results              []result
	requiredLetterCounts map[rune]int
	prohibitedLetters1   map[rune]struct{}
	prohibitedLetters2   map[rune]struct{}
	prohibitedLetters3   map[rune]struct{}
	prohibitedLetters4   map[rune]struct{}
	prohibitedLetters5   map[rune]struct{}
}

func newHistory() *history {
	h := history{
		requiredLetterCounts: make(map[rune]int, 5),
		prohibitedLetters1:   make(map[rune]struct{}, 26),
		prohibitedLetters2:   make(map[rune]struct{}, 26),
		prohibitedLetters3:   make(map[rune]struct{}, 26),
		prohibitedLetters4:   make(map[rune]struct{}, 26),
		prohibitedLetters5:   make(map[rune]struct{}, 26),
	}
	return &h
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
	requiredLetterCounts := make(map[rune]int, 5)
	for i, si := range r.score {
		gi := rune(r.guess[i])
		p := h.prohibitedLetters(i)
		switch si {
		case 'c':
			if _, ok := p[gi]; ok {
				return fmt.Errorf("%c was prohibited at inde %v, but is now supposedly correct", si, i)
			}
			for l := 'a'; l <= 'z'; l++ {
				if l != gi {
					p[l] = struct{}{}
				}
			}
			requiredLetterCounts[gi]++
		case 'a':
			p[gi] = struct{}{}
			if len(p) == 26 {
				return fmt.Errorf("all letters prohibited at index %v", i)
			}
			requiredLetterCounts[gi]++
		case 'n':
			if _, ok := h.requiredLetterCounts[gi]; ok {
				return fmt.Errorf("%c was previously required to be in word, but is prohibited", gi)
			}
			for j := range r.score {
				pj := h.prohibitedLetters(j)
				pj[gi] = struct{}{}
				if len(pj) == 26 {
					return fmt.Errorf("all letters prohibited at index %v", i)
				}
			}
		}
	}
	if err := h.mergeRequiredLetterCounts(requiredLetterCounts); err != nil {
		return fmt.Errorf("merging required letters: %v", err)
	}
	return nil
}

func (h *history) mergeRequiredLetterCounts(extraLetterCounts map[rune]int) error {
	c := 0
	for ch, n := range h.requiredLetterCounts {
		if n2, ok := extraLetterCounts[ch]; ok {
			if n2 > n {
				extraLetterCounts[ch] = n2 - n
			} else {
				delete(extraLetterCounts, ch)
			}
		}
		c += n
	}
	for ch, n := range extraLetterCounts {
		if c+n > 5 {
			return fmt.Errorf("more than five letters are now required")
		}
		h.requiredLetterCounts[ch] += n
	}
	return nil
}

func (h *history) prohibitedLetters(index int) map[rune]struct{} {
	switch index {
	case 0:
		return h.prohibitedLetters1
	case 1:
		return h.prohibitedLetters2
	case 2:
		return h.prohibitedLetters3
	case 3:
		return h.prohibitedLetters4
	case 4:
		return h.prohibitedLetters5
	default:
		panic(fmt.Errorf("unknown prohibited letter index: %v", index))
	}
}

func (h *history) allows(w string) bool {
	letterCounts := make(map[rune]int, 5)
	for i, ch := range w {
		p := h.prohibitedLetters(i)
		if _, ok := p[ch]; ok {
			return false
		}
		letterCounts[ch]++
	}
	for ch, n := range h.requiredLetterCounts {
		if n2, ok := letterCounts[ch]; !ok || n < n2 {
			return false
		}
	}
	return true
}

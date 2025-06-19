package guess

import (
	"fmt"
	"io"
	"strings"

	words "github.com/jacobpatterson1549/wordle-cheater"
)

const numLetters = 5

// Guess is a word that might be the answer
type Guess string

// New reads the next word from the reader.  It may be invalid.
func New(word string) Guess {
	word = strings.ToLower(word)
	g := Guess(word)
	return g
}

// Scan prompts for a guess on the ReadWriter until a valid one is given or an io error occurs
func Scan(rw io.ReadWriter, m words.Words) (*Guess, error) {
	for {
		fmt.Fprintf(rw, "Enter guess (%v letters): ", numLetters)
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		g := New(word)
		if err := g.Validate(m); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &g, nil
	}
}

// Validate ensures the guess is <<numLetters>> letters long and is in the words list (if a list is provided)
func (g Guess) Validate(m words.Words) error {
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

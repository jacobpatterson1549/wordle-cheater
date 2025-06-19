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

// New prompts for a guess on the ReadWriter until a valid one is given or an io error occurs
func New(rw io.ReadWriter, m words.Words) (*Guess, error) {
	for {
		fmt.Fprintf(rw, "Enter guess (%v letters): ", numLetters)
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		word = strings.ToLower(word)
		g := Guess(word)
		if err := g.validate(m); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &g, nil
	}
}

// validate ensures the guess is <<numLetters>> letters long and is in the words list (if a list is provided)
func (g Guess) validate(m words.Words) error {
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

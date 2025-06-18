package main

import (
	"fmt"
	"io"
	"strings"
)

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

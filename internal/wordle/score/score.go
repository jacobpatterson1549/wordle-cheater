package score

import (
	"fmt"
	"io"
	"strings"
)

// Score is a <<numLetters>>-letter string made up of {c,a,n}.
// * The letter c indicates that a letter from a guess is in the correct position.
// * The letter a indicates that a letter from a guess is in the answer, but in a different position.
// * The letter n indicates that a letter from a guess is not anywhere in the answer.
type Score string

const AllCorrect Score = "ccccc"

// New prompts for a score on the ReadWriter until a valid one is given or an io error occurs
func New(rw io.ReadWriter) (*Score, error) {
	for {
		fmt.Fprintf(rw, "Enter score: ")
		var word string
		if _, err := fmt.Fscan(rw, &word); err != nil {
			return nil, fmt.Errorf("scanning guess: %v", err)
		}
		word = strings.ToLower(word)
		s := Score(word)
		if err := s.validate(); err != nil {
			fmt.Fprintf(rw, "%v\n", err)
			continue
		}
		return &s, nil
	}
}

// validate ensures the score is <<numLetters>> letters long and consists only of the {c,a,n} letters
func (s Score) validate() error {
	const numLetters = 5
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

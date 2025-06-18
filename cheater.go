package main

import (
	"fmt"
	"io"
)

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
		if *s == allCorrect {
			return nil
		}

		r := result{
			guess: *g,
			score: *s,
		}
		h.addResult(r, availableWords)

		if err := availableWords.scanShowPossible(rw); err != nil {
			return err
		}
	}
}

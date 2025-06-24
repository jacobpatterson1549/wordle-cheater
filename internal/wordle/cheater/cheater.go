package cheater

import (
	"fmt"
	"io"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/guess"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/result"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/score"
)

func RunWordleCheater(rw io.ReadWriter, wordsText string) error {
	const numLetters = 5
	allWords, err := words.New(wordsText)
	if err != nil {
		return fmt.Errorf("loading words: %v", err)
	}
	availableWords := allWords.Copy()

	fmt.Fprintf(rw, "Running wordle-cheater\n")
	fmt.Fprintf(rw, " * Guesses and scores are %v letters long\n", numLetters)
	fmt.Fprintf(rw, " * Scores are only made of the following letters:\n")
	fmt.Fprintf(rw, "   C - if a letter is in the word and in the correct location\n")
	fmt.Fprintf(rw, "   A - if a letter is in the word, but in the wrong location\n")
	fmt.Fprintf(rw, "   N - if a letter is not in the word\n")
	fmt.Fprintf(rw, "The app runs until the correct word is found from a guess with only correct letters.\n\n")

	var h result.History
	for {
		g, err := guess.Scan(rw, *allWords)
		if err != nil {
			return err
		}

		s, err := score.Scan(rw)
		if err != nil {
			return err
		}
		if *s == score.AllCorrect {
			return nil
		}

		r := result.Result{
			Guess: *g,
			Score: *s,
		}
		h.AddResult(r, availableWords)

		if err := availableWords.ScanShowPossible(rw); err != nil {
			return err
		}
	}
}

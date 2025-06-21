package main

import (
	"fmt"
	"io"
	"os"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/spelling_bee"
)

func main() {
	runSpellingBee(os.Stdin, os.Stdout)

}

func runSpellingBee(r io.Reader, w io.Writer) {
	var sb spelling_bee.SpellingBee
	sb.MinLength = 4

	fmt.Fprint(w, "enter central letter: ")
	fmt.Fscan(r, &sb.OtherLetters)
	if len(sb.OtherLetters) != 1 {
		fmt.Fprintln(w, "expected 1 central letter")
		return
	}
	sb.CentralLetter = rune(sb.OtherLetters[0])

	fmt.Print("enter other letters: ")
	fmt.Fscan(r, &sb.OtherLetters)

	if len(sb.OtherLetters) != 6 {
		fmt.Fprintln(w, "expected 6 other letters")
		return
	}

	fmt.Println("available words: (score first)")
	words := sb.Words(words.WordsTextFile)
	for _, v := range words {
		fmt.Fprint(w, v.Score, " ", v.Value)
		if v.IsPangram {
			fmt.Print(w, " (PANGRAM!)")
		}
		fmt.Fprintln(w)
	}
}

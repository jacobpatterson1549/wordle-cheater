// Package main runs a command-line-interface program to help cheat in the popular Wordle game
package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
)

//go:embed build/words.txt
var wordsTextFile string

// numLetters is the length of the words
const numLetters = 5

// init ensures the program is set up properly
func init() {
	if err := allCorrect.validate(); err != nil {
		panic(fmt.Errorf("all correct string is not valid: %v", err))
	}
}

// main runs wordle-cheater on the command-line using stdin and stdout
func main() {
	rw := struct {
		io.Reader
		io.Writer
	}{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	if err := runWordleCheater(rw, wordsTextFile); err != nil {
		panic(fmt.Errorf("running wordle: %v", err))
	}
}

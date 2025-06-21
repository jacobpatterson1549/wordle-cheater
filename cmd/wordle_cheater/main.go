// Package main runs a command-line-interface program to help cheat in the popular Wordle game
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/cheater"
)

// main runs wordle-cheater on the command-line using stdin and stdout
func main() {
	rw := struct {
		io.Reader
		io.Writer
	}{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	if err := cheater.RunWordleCheater(rw); err != nil {
		panic(fmt.Errorf("running wordle: %v", err))
	}
}

# wordle-cheater

[![Docker Image CI](https://github.com/jacobpatterson1549/wordle-cheater/actions/workflows/docker-image.yml/badge.svg)](https://github.com/jacobpatterson1549/wordle-cheater/actions/workflows/docker-image.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jacobpatterson1549/wordle-cheater)](https://goreportcard.com/report/github.com/jacobpatterson1549/wordle-cheater)
[![GoDoc](https://godoc.org/github.com/jacobpatterson1549/wordle-cheater?status.svg)](https://godoc.org/github.com/jacobpatterson1549/wordle-cheater)


A command-line application that can help users play [Wordle](https://www.powerlanguage.co.uk/wordle).

The application shows possible words after combining the results of previous guesses/scores that are entered by the user.
It runs until a correct guess, with a score of "ccccc" is entered.

## usage

[Go 1.25](https://golang.org/dl/) is used to build the application.
* Version 1.17 is needed to embed the words list in the executable file.
* Version 1.24 is for recent updates: security, range expressions.
* Version 1.25 is for simplified documentation generation.

[Aspell](https://github.com/GNUAspell/aspell) is used to generate the words list.

[Make](https://www.gnu.org/software/make/) is used to automate code compilation.  The command `make` builds the application into an executable file.

The command `make` builds the application in a terminal.  This generates the word list, tests the code, and compiles the http server and command-line-interfaces.  The programs are placed in the build/bin folder.  The application runs until the correct guess is entered or control-c is pressed.

To build the application to be run on other operating systems/architectures, set the GO_ARGS flag when running `make`.  An example of this is `make build/bin/wordle_cheater GO_ARGS="GOOS=windows GOARCH=amd64" OBJ="wordle-cheater.exe"`.  This builds `build/bin/wordle_cheater.exe`, a version of the application that runs on 64-bit versions of Windows.  To list available architectures, run `go tool dist list` to display GOOS/GOARCH combinations.

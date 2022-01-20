# wordle-cheater

A command-line application that can help users play [Wordle](https://www.powerlanguage.co.uk/wordle).

The application shows possible words after combining the results of previous guesses/scores that are entered by the user.
It runs until a correct guess, with a score of "ccccc" is entered.

## usage

[Go 1.17](https://golang.org/dl/) is used to build the application. Version 1.17 is needed to embed the words list in the executable file.

[Aspell](https://github.com/GNUAspell/aspell) is used to generate the words list.

[Make](https://www.gnu.org/software/make/) is used to automate code compilation.  The command `make` builds the application into an executable file.

The command `make run` executes the application in a terminal.  This generates the word list, tests the code, and runs the command-line-interface application in the terminal.  The application runs until the correct guess is entered or control-c is pressed.

To build the application to be run on other operating systems/architectures, set the GO_ARGS flag when running `make`.  An examlpe of this is `make GO_ARGS="GOOS=windows GOARCH=amd64" OBJ="wordle-cheater.exe"`.  This builds `build/wordle-cheater.exe`, a version of the application that runs on 64-bit versions of Windows.  To list available architectures, run `go tool dist list` to display GOOS/GOARCH combinations.
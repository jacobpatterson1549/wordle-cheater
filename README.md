# wordle-cheater

A utility that can help users play [Wordle](https://www.powerlanguage.co.uk/wordle).

## usage

[Go 1.17](https://golang.org/dl/) is used to build the application. 1.17 is needed to embed the words list in the executable file.

[Aspell](https://github.com/GNUAspell/aspell) is used to generate the words list

[Make](https://www.gnu.org/software/make/) is used to help automate code compilation

To run the code, run `make run` in a terminal.  This generates the word list, tests the code, and runs the command-line-interface program in the terminal.  The application runs until the correct guess is entered or control-c is pressed.
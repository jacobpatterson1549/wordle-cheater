package cheater

import (
	"bufio"
	"strings"
	"testing"

	words "github.com/jacobpatterson1549/wordle-cheater"
)

func TestRunWordleCheater(t *testing.T) {
	o := words.WordsTextFile
	defer func() { words.WordsTextFile = o }()
	tests := []struct {
		readTokens string
		wordsText  string
		wantErr    bool
	}{
		{
			readTokens: "smart ccccc",
			wordsText:  "smart",
		},
		{
			readTokens: "smart ccccc",
			wordsText:  "Smart",
			wantErr:    true,
		},
		{
			readTokens: "dummy nnnnn n smart ccccc",
			wordsText:  "smart",
		},
		{
			wordsText: "some tiny text",
			wantErr:   true, // words too short
		},
		{
			wantErr: true, // EOF guess
		},
		{
			readTokens: "guess",
			wantErr:    true, // EOF score
		},
		{
			readTokens: "guess nnnnn",
			wantErr:    true, // EOF scanShowPossible
		},
		{
			readTokens: "apple ncccc n berry ncccc",
			wantErr:    true, // too many required letters
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.readTokens)),
			Writer: bufio.NewWriter(&buf),
		}
		words.WordsTextFile = test.wordsText
		gotErr := RunWordleCheater(rw)
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error running wordle cheater", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error running wordle cheater", i)
		}
	}
}

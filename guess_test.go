package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestNewGuess(t *testing.T) {
	simpleOut := "Enter guess (5 letters): "
	tests := []struct {
		in       string
		wantOut  string
		allWords map[string]struct{}
		want     guess
		wantErr  bool
	}{
		{
			wantErr: true, // input EOF
		},
		{
			in:      "happy",
			wantOut: simpleOut,
			want:    "happy",
		},
		{
			in:       "happy",
			wantOut:  simpleOut,
			allWords: map[string]struct{}{"happy": {}},
			want:     "happy",
		},
		{
			in:       "HAPPY",
			wantOut:  simpleOut,
			allWords: map[string]struct{}{"happy": {}},
			want:     "happy",
		},
		{
			in:       "tiny error happy",
			wantOut:  "Enter guess (5 letters): guess must be 5 letters long\nEnter guess (5 letters): error is not a word\nEnter guess (5 letters): ",
			allWords: map[string]struct{}{"happy": {}},
			want:     "happy",
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.in)),
			Writer: bufio.NewWriter(&buf),
		}
		got, gotErr := newGuess(rw, test.allWords)
		rw.Flush()
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error creating guess", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error creating guess: %v", i, gotErr)
		case test.want != *got:
			t.Errorf("test %v: guesses not equal:\nwanted: %v\ngot:    %v", i, test.want, *got)
		case test.wantOut != buf.String():
			t.Errorf("test %v: outputs not equal creating guess:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())
		}
	}
}

package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNewWords(t *testing.T) {
	tests := []struct {
		input   string
		want    *words
		wantErr bool
	}{
		{
			wantErr: true, // too short (length 0)
		},
		{
			input:   "tiny", // too short
			wantErr: true,
		},
		{
			input: "apple\nberry",
			want:  &words{"apple": {}, "berry": {}},
		},
		{
			input:   "APPLE", // uppercase
			wantErr: true,
		},
	}
	for i, test := range tests {
		got, err := newWords(test.input)
		switch {
		case test.wantErr:
			if err == nil {
				t.Errorf("test %v: wanted error", i)
			}
		case err != nil:
			t.Errorf("test %v: unwanted error: %v", i, err)
		case !reflect.DeepEqual(test.want, got):
			t.Errorf("test %v: words not equal:\nwanted: %v\ngot:    %v", i, test.want, got)
		}
	}
}

func TestNewGuess(t *testing.T) {
	tests := []struct {
		in       string
		wantOut  string
		allWords map[string]struct{}
		want     guess
		wantErr  bool
	}{
		{
			in:      "happy",
			wantOut: "Enter guess (five letters): ",
			want:    "happy",
		},
		{
			in:       "happy",
			wantOut:  "Enter guess (five letters): ",
			allWords: map[string]struct{}{"happy": {}},
			want:     "happy",
		},
		{
			in:       "HAPPY",
			wantOut:  "Enter guess (five letters): ",
			allWords: map[string]struct{}{"happy": {}},
			want:     "happy",
		},
		{
			in:       "tiny error happy",
			wantOut:  "Enter guess (five letters): guess must be 5 letters long\nEnter guess (five letters): error is not a word\nEnter guess (five letters): ",
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
				t.Errorf("test %v: wanted error", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error: %v", i, gotErr)
		case test.want != *got:
			t.Errorf("test %v: guesses not equal:\nwanted: %v\ngot:    %v", i, test.want, *got)
		case test.wantOut != buf.String():
			t.Errorf("test %v: outputs not equal:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())

		}
	}
}

func TestNewScore(t *testing.T) {
	tests := []struct {
		in      string
		wantOut string
		want    score
		wantErr bool
	}{
		{
			in:      "ccccc",
			wantOut: "Enter score: ",
			want:    "ccccc",
		},
		{
			in:      "cCcCc",
			wantOut: "Enter score: ",
			want:    "ccccc",
		},
		{
			in:      "nac apple canac",
			wantOut: "Enter score: score must be 5 letters long\nEnter score: must be only the following letters: C, A, N\nEnter score: ",
			want:    "canac",
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.in)),
			Writer: bufio.NewWriter(&buf),
		}
		got, gotErr := newScore(rw)
		rw.Flush()
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error: %v", i, gotErr)
		case test.want != *got:
			t.Errorf("test %v: scores not equal:\nwanted: %v\ngot:    %v", i, test.want, *got)
		case test.wantOut != buf.String():
			t.Errorf("test %v: outputs not equal:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())

		}
	}
}

func TestScoreAllCorrect(t *testing.T) {
	tests := []struct {
		score
		want bool
	}{
		{
			score: "cccca",
		},
		{
			score: "c",
		},
		{
			score: "ccccc",
			want:  true,
		},
	}
	for i, test := range tests {
		if want, got := test.want, test.score.allCorrect(); want != got {
			t.Errorf("test %v: allCorrect values not equal for %q: wanted %v, got %v", i, test.score, want, got)
		}
	}
}

func TestResultValidate(t *testing.T) {
	tests := []struct {
		result
		wantErr bool
	}{
		{
			result: result{
				guess: "tiny",
				score: "ccccc",
			},
			wantErr: true,
		},
		{
			result: result{
				guess: "large",
				score: "n",
			},
			wantErr: true,
		},
		{
			result: result{
				guess: "happy",
				score: "cannc",
			},
		},
	}
	for i, test := range tests {
		gotErr := test.result.validate()
		if want, got := test.wantErr, gotErr != nil; want != got {
			t.Errorf("test %v: validated values not equal: wanted error: %v, got error: %v (%v)", i, want, got, gotErr)
		}
	}
}

func TestNewHistory(t *testing.T) {
	h := newHistory()
	switch {
	case h.requiredLetterCounts == nil:
		t.Errorf("requiredLetterCounts is nil")
	case h.prohibitedLetters1 == nil:
		t.Errorf("prohibitedLetters1 is nil")
	case h.prohibitedLetters2 == nil:
		t.Errorf("prohibitedLetters2 is nil")
	case h.prohibitedLetters3 == nil:
		t.Errorf("prohibitedLetters3 is nil")
	case h.prohibitedLetters4 == nil:
		t.Errorf("prohibitedLetters4 is nil")
	case h.prohibitedLetters5 == nil:
		t.Errorf("prohibitedLetters5 is nil")
	}
}

func TestAddResult(t *testing.T) {
	s := []string{"nasty", "alley", "early", "great", "ready", "touch"}
	allWords := make(words, len(s))
	for _, w := range s {
		allWords[w] = struct{}{}
	}
	r := result{
		guess: "nasty",
		score: "nannc",
	}
	want := &history{
		results:              []result{r},
		requiredLetterCounts: map[rune]int{'a': 1, 'y': 1},
		prohibitedLetters1:   map[rune]struct{}{'n': {}, 's': {}, 't': {}},
		prohibitedLetters2:   map[rune]struct{}{'n': {}, 's': {}, 't': {}, 'a': {}},
		prohibitedLetters3:   map[rune]struct{}{'n': {}, 's': {}, 't': {}},
		prohibitedLetters4:   map[rune]struct{}{'n': {}, 's': {}, 't': {}},
		prohibitedLetters5:   map[rune]struct{}{'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {}, 'n': {}, 'o': {}, 'p': {}, 'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {}, 'z': {}},
	}
	wantWords := words{
		"alley": {},
		"ready": {},
	}
	h := newHistory()
	h.addResult(r, &allWords)
	got := h
	gotWords := allWords
	switch {
	case !reflect.DeepEqual(want, got):
		t.Errorf("histories not equal:\nwanted: %+v\ngot:    %+v", want, got)
	case !reflect.DeepEqual(wantWords, gotWords):
		t.Errorf("words not equal after result added to history:\nwanted: %+v\ngot:    %+v", wantWords, gotWords)
	}
}

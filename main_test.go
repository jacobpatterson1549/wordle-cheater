package main

import (
	"bufio"
	"fmt"
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
			want: &words{}, // no words
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

func TestWordsSorted(t *testing.T) {
	words := words{
		"abbey": {},
		"weary": {},
		"gravy": {},
	}
	if want, got := "abbey,gravy,weary", words.sorted(); want != got {
		t.Errorf("sorted words not equal:\nwanted: %q\ngot:    %q", want, got)
	}
}

func TestWordsCopy(t *testing.T) {
	w := "magic"
	a := words{
		w: {},
	}
	b := a.copy()
	if !reflect.DeepEqual(a, *b) {
		t.Errorf("copied values should be equal:\nwanted: %v\ngot:    %v", a, b)
	}
	delete(a, w)
	if reflect.DeepEqual(a, *b) {
		t.Errorf("copy should be to a different location in memory, both were %p", &a)
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

func TestHistoryAddResult(t *testing.T) {
	s := []string{"nasty", "alley", "early", "great", "ready", "touch"}
	allWords := make(words, len(s))
	for _, w := range s {
		allWords[w] = struct{}{}
	}
	r := result{
		guess: "nasty",
		score: "nannc",
	}
	want := history{
		requiredLetters:    []rune{'a', 'y'},
		prohibitedLetters1: newCharSetHelper(t, 'n', 's', 't'),
		prohibitedLetters2: newCharSetHelper(t, 'n', 's', 't', 'a'),
		prohibitedLetters3: newCharSetHelper(t, 'n', 's', 't'),
		prohibitedLetters4: newCharSetHelper(t, 'n', 's', 't'),
		prohibitedLetters5: newCharSetHelper(t, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'z'),
	}
	wantWords := words{
		"alley": {},
		"ready": {},
	}
	var h history
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

func TestHistoryMergeRequiredLetters(t *testing.T) {
	tests := []struct {
		history
		newRequiredLetters []rune
		want               history
		wantErr            bool
	}{
		{},
		{
			newRequiredLetters: []rune{'a', 'b', 'c'},
			want:               history{requiredLetters: []rune{'a', 'b', 'c'}},
		},
		{
			history:            history{requiredLetters: []rune{'a', 'b'}},
			newRequiredLetters: []rune{'a', 'a'},
			want:               history{requiredLetters: []rune{'a', 'b', 'a'}},
		},
		{
			history:            history{requiredLetters: []rune{'a', 'a', 'a'}},
			newRequiredLetters: []rune{'a', 'a', 'b', 'b', 'c'},
			wantErr:            true,
		},
	}
	for i, test := range tests {
		got := test.history
		err := got.mergeRequiredLetters(test.newRequiredLetters...)
		switch {
		case test.wantErr:
			if err == nil {
				t.Errorf("test %v: wanted error", i)
			}
		case err != nil:
			t.Errorf("test %v: unwanted error: %v", i, err)
		case !reflect.DeepEqual(test.want, got):
			t.Errorf("test %v histories not equal:\nwanted: %v\ngot:    %v", i, test.want, got)
		}
	}
}

func TestCharSet(t *testing.T) {
	for ch := rune('a'); ch <= 'z'; ch++ {
		var cs charSet
		if cs.has(ch) {
			t.Errorf("%c in charSet before it is added", ch)
		}
		cs.add(ch)

		if !cs.has(ch) {
			t.Errorf("%c not in charSet after it is added", ch)
		}
	}
	t.Run("isFull", func(t *testing.T) {
		var cs charSet
		for ch := rune('a'); ch <= 'z'; ch++ {
			if cs.isFull() {
				t.Fatalf("charSet is full before adding %v", ch)
			}
			cs.add(ch)
		}
		if !cs.isFull() {
			t.Fatalf("wanted charSet to be is full after adding a-z")
		}
	})
	t.Run("Stringer", func(t *testing.T) {
		cs := newCharSetHelper(t, 'f', 'y', 'r', 'o', 't')
		if want, got := "[forty]", cs.String(); want != got {
			t.Errorf("wanted %q, got %q", want, got)
		}
	})
	badChars := []rune{'?', 'A', 'Z', ' ', '!', '`', '{', '\n', 0, 0x7F, 0xFF}
	for _, ch := range badChars {
		t.Run(fmt.Sprintf("bad-add-0x%x", ch), func(t *testing.T) {
			var cs charSet
			if cs.has(ch) {
				t.Errorf("bad character 0x%x in charSet", ch)
			}
			defer func() {
				r := recover()
				if _, ok := r.(error); r == nil || !ok {
					t.Errorf("expected panic error adding bad character")
				}
			}()
			cs.add(ch)
		})
	}
}

func newCharSetHelper(t *testing.T, chars ...rune) charSet {
	t.Helper()
	var cs charSet
	for _, ch := range chars {
		cs.add(ch)
	}
	return cs
}

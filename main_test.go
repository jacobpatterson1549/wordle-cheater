package main

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestRunWordleCheater(t *testing.T) {
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
			wantErr:    true, // to mainy required letters
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.readTokens)),
			Writer: bufio.NewWriter(&buf),
		}
		gotErr := runWordleCheater(rw, test.wordsText)
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error running wordle cheater", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error runnig wordle cheater", i)
		}
	}
}

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
		got, gotErr := newWords(test.input)
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error creating words", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error creating words: %v", i, gotErr)
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

func TestWordsScanShowPossible(t *testing.T) {
	tests := []struct {
		words
		in      string
		wantOut string
		wantErr bool
	}{
		{
			wantErr: true, // input EOF
		},
		{
			in:      "n",
			wantOut: "show possible words [Yn]: ",
		},
		{
			words:   words{"apple": {}, "berry": {}, "cakes": {}},
			in:      "NAH",
			wantOut: "show possible words [Yn]: ",
		},
		{
			words:   words{"apple": {}, "berry": {}, "cakes": {}},
			in:      "yes",
			wantOut: "show possible words [Yn]: remaining valid words: apple,berry,cakes\n",
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.in)),
			Writer: bufio.NewWriter(&buf),
		}
		gotErr := test.words.scanShowPossible(rw)
		rw.Flush()
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error scanning to show possible words", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error scanning to show possible words: %v", i, gotErr)
		case test.wantOut != buf.String():
			t.Errorf("test %v: outputs not equal for possible words:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())
		}
	}
}

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

func TestNewScore(t *testing.T) {
	tests := []struct {
		in      string
		wantOut string
		want    score
		wantErr bool
	}{
		{
			wantErr: true, // input EOF
		},
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
				t.Errorf("test %v: wanted error creating score", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error creating score: %v", i, gotErr)
		case test.want != *got:
			t.Errorf("test %v: scores not equal:\nwanted: %v\ngot:    %v", i, test.want, *got)
		case test.wantOut != buf.String():
			t.Errorf("test %v: new score outputs not equal:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())
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

func TestHistoryAddResultInvalid(t *testing.T) {
	var h history
	var r result
	if gotErr := h.addResult(r, nil); gotErr == nil {
		t.Errorf("wanted error adding invalid result")
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
		requiredLetters: []rune{'a', 'y'},
		prohibitedLetters: [numLetters]charSet{
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't', 'a'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'z'),
		},
	}
	wantWords := words{
		"alley": {},
		"ready": {},
	}
	var h history
	gotErr := h.addResult(r, &allWords)
	got := h
	gotWords := allWords
	switch {
	case gotErr != nil:
		t.Errorf("unwanted error: %v", gotErr)
	case !reflect.DeepEqual(want, got):
		t.Errorf("histories not equal:\nwanted: %+v\ngot:    %+v", want, got)
	case !reflect.DeepEqual(wantWords, gotWords):
		t.Errorf("words not equal after result added to history:\nwanted: %+v\ngot:    %+v", wantWords, gotWords)
	}
}

func TestHistoryMergeResult(t *testing.T) {
	t.Run("invalid-merges", func(t *testing.T) {
		tests := []struct {
			history
			result
		}{
			{
				history: history{
					prohibitedLetters: [numLetters]charSet{
						0: newCharSetHelper(t, 'a'),
					},
				},
				result: result{
					guess: "apple",
					score: "cnnna", // 'a' was prohibited as the first letter
				},
			},
			{
				history: history{
					prohibitedLetters: [numLetters]charSet{
						0: newCharSetHelper(t, 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'),
					},
				},
				result: result{
					guess: "alter",
					score: "annnn", // all letters prohibited for first letter
				},
			},
			{
				history: history{
					requiredLetters: []rune{'a'},
				},
				result: result{
					guess: "apple",
					score: "nannc", // 'a' was required
				},
			},
			{
				history: history{
					prohibitedLetters: [numLetters]charSet{
						0: newCharSetHelper(t, 'a', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'),
					},
				},
				result: result{
					guess: "berry",
					score: "nnnnn", // all letters prohibited for first letter
				},
			},
			{
				history: history{
					requiredLetters: []rune{'a', 'b', 'c', 'd'},
				},
				result: result{
					guess: "apple",
					score: "cnnac", // too many letters are now required
				},
			},
		}
		for i, test := range tests {
			if gotErr := test.history.mergeResult(test.result); gotErr == nil {
				t.Errorf("test %v: wanted error for invalid merge", i)
			}
		}
	})
}

func TestHistoryHasRequiredLetter(t *testing.T) {
	tests := []struct {
		history
		letter          rune
		newScoreLetters []rune
		want            bool
	}{
		{
			letter:  'a',
			history: history{requiredLetters: []rune{'b', 'a', 't'}},
			want:    true,
		},
		{
			letter:          'a',
			history:         history{requiredLetters: []rune{'b', 'a', 't'}},
			newScoreLetters: []rune{'z'},
			want:            true,
		},
		{
			letter:          'o',
			history:         history{requiredLetters: []rune{'a', 't'}},
			newScoreLetters: []rune{'z', 'o', 'o'},
			want:            true,
		},
		{
			letter:          'r',
			history:         history{requiredLetters: []rune{'t'}},
			newScoreLetters: []rune{'x'},
			want:            false,
		},
	}
	for i, test := range tests {
		if want, got := test.want, test.history.hasRequiredLetter(test.letter, test.newScoreLetters...); want != got {
			t.Errorf("test %v: wanted %v, got %v", i, want, got)
		}
	}
}

func TestHistoryMergeRequiredLetters(t *testing.T) {
	tests := []struct {
		history
		newRequiredLetters []rune
		want               history
		wantErr            bool
	}{
		{
			want: history{requiredLetters: []rune{}},
		},
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
		gotErr := got.mergeRequiredLetters(test.newRequiredLetters...)
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error merging required letters", i)
			}
			if len(got.requiredLetters) > numLetters {
				t.Errorf("test %v: wanted required letters not to be modified to invalid state", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error merging required letters: %v", i, gotErr)
		case !reflect.DeepEqual(test.want, got):
			t.Errorf("test %v histories not equal:\nwanted: %v\ngot:    %v", i, test.want, got)
		}
	}
}

func TestHistoryAllowsWord(t *testing.T) {
	tests := []struct {
		word string
		want bool
	}{
		{"batty", true},
		{"fatty", false},
		{"party", false},
	}
	for i, test := range tests {
		h := history{
			requiredLetters: []rune{'t', 'a', 't'},
			prohibitedLetters: [numLetters]charSet{
				0: newCharSetHelper(t, 'f'),
			},
		}
		if want, got := test.want, h.allows(test.word); want != got {
			t.Errorf("test %v: wanted %v, got %v", i, want, got)
		}
	}
}

func TestCharSetHas(t *testing.T) {
	for ch := rune('a'); ch <= 'z'; ch++ {
		var cs charSet
		if cs.Has(ch) {
			t.Errorf("%c in charSet before it is added", ch)
		}
		cs.Add(ch)
		if !cs.Has(ch) {
			t.Errorf("%c not in charSet after it is added", ch)
		}
	}
}

func TestCharSetIsFull(t *testing.T) {
	var cs charSet
	for ch := rune('a'); ch <= 'z'; ch++ {
		if cs.IsFull() {
			t.Fatalf("charSet is full before adding %v", ch)
		}
		cs.Add(ch)
	}
	if !cs.IsFull() {
		t.Fatalf("wanted charSet to be is full after adding a-z")
	}
}

func TestCharSetString(t *testing.T) {
	cs := newCharSetHelper(t, 'f', 'y', 'r', 'o', 't')
	if want, got := "[forty]", cs.String(); want != got {
		t.Errorf("wanted %q, got %q", want, got)
	}
}

func TestCharSetBadChars(t *testing.T) {
	badChars := []rune{'?', 'A', 'Z', ' ', '!', '`', '{', '\n', 0, 0x7F, 0xFF}
	for _, ch := range badChars {
		t.Run(fmt.Sprintf("bad-add-0x%x", ch), func(t *testing.T) {
			var cs charSet
			if cs.Has(ch) {
				t.Errorf("bad character 0x%x in charSet", ch)
			}
			defer func() {
				r := recover()
				if _, ok := r.(error); r == nil || !ok {
					t.Errorf("expected panic error adding bad character")
				}
			}()
			cs.Add(ch)
		})
	}
}

func newCharSetHelper(t *testing.T, chars ...rune) charSet {
	t.Helper()
	var cs charSet
	for _, ch := range chars {
		cs.Add(ch)
	}
	return cs
}

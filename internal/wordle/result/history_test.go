package result

import (
	"reflect"
	"testing"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/char_set"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/guess"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/score"
	"github.com/jacobpatterson1549/wordle-cheater/internal/words"
)

func TestHistoryAddResult(t *testing.T) {
	s := []string{"nasty", "alley", "early", "great", "ready", "touch"}
	allWords := make(words.Words, len(s))
	for _, w := range s {
		allWords[w] = struct{}{}
	}
	r := Result{
		Guess: "nasty",
		Score: "nannc",
	}
	want := History{
		correctLetters: [numLetters]rune{
			4: 'y',
		},
		almostLetters: []rune{'a', 'y'},
		prohibitedLetters: [numLetters]char_set.CharSet{
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't', 'a'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't'),
		},
	}
	wantWords := words.Words{
		"alley": {},
		"ready": {},
	}
	var h History
	got := h
	got.AddResult(r, &allWords)
	gotWords := allWords
	switch {
	case !reflect.DeepEqual(want, got):
		t.Errorf("histories not equal:\nwanted: %+v\ngot:    %+v", want, got)
	case !reflect.DeepEqual(wantWords, gotWords):
		t.Errorf("words not equal after result added to history:\nwanted: %+v\ngot:    %+v", wantWords, gotWords)
	}
}

func TestHistoryMergeResult(t *testing.T) {
	tests := []struct {
		History
		guess.Guess
		score.Score
		want History
	}{
		{
			Guess: "treat",
			Score: "nannc",
			want: History{
				correctLetters: [numLetters]rune{
					4: 't',
				},
				almostLetters: []rune{'r', 't'},
				prohibitedLetters: [numLetters]char_set.CharSet{
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a', 'r'),
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a'),
				},
			},
		},
		{
			Guess: "shove",
			Score: "accnc",
			want: History{
				correctLetters: [numLetters]rune{
					1: 'h',
					2: 'o',
					4: 'e',
				},
				almostLetters: []rune{'s', 'h', 'o', 'e'},
				prohibitedLetters: [numLetters]char_set.CharSet{
					newCharSetHelper(t, 's', 'v'),
					newCharSetHelper(t, 'v'),
					newCharSetHelper(t, 'v'),
					newCharSetHelper(t, 'v'),
					newCharSetHelper(t, 'v'),
				},
			},
		},
	}
	for i, test := range tests {
		r := Result{
			Guess: test.Guess,
			Score: test.Score,
		}
		want := test.want
		got := test.History
		got.mergeResult(r)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("test %v (for guess %q): histories not equal:\nwanted: %+v\ngot:    %+v", i, test.Guess, want, got)
		}
	}
}

func TestString(t *testing.T) {
	h := History{
		correctLetters: [numLetters]rune{
			4: 'q',
		},
		almostLetters: []rune{'c', 'a', 'b'},
		prohibitedLetters: [numLetters]char_set.CharSet{
			1: newCharSetHelper(t, 'z', 'e', 'r'),
			2: newCharSetHelper(t, 'z', 'x', 'a'),
		},
	}
	want := `{correctLetters:????q almostLetters:[c a b] prohibitedLetters:[[] [erz] [axz] [] []]}`
	got := h.String()
	if want != got {
		t.Errorf("history Strings not equal:\nwanted: %+v\ngot:    %+v", want, got)
	}
}

func TestHistoryMergeRequiredLetters(t *testing.T) {
	tests := []struct {
		History
		newRequiredLetters []rune
		want               History
	}{
		{
			want: History{almostLetters: []rune{}},
		},
		{
			newRequiredLetters: []rune{'a', 'b', 'c'},
			want:               History{almostLetters: []rune{'a', 'b', 'c'}},
		},
		{
			History:            History{almostLetters: []rune{'a', 'b'}},
			newRequiredLetters: []rune{'a', 'a'},
			want:               History{almostLetters: []rune{'a', 'b', 'a'}},
		},
		{
			History:            History{almostLetters: []rune{'a', 'a', 'a'}},
			newRequiredLetters: []rune{'a', 'a', 'b', 'b', 'c'},
			want:               History{almostLetters: []rune{'a', 'a', 'a', 'b', 'b', 'c'}}, // this will prohibit all words
		},
	}
	for i, test := range tests {
		got := test.History
		got.mergeRequiredLetters(test.newRequiredLetters)
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("test %v histories not equal:\nwanted: %v\ngot:    %v", i, test.want, got)
		}
	}
}

func TestHistoryAllowsWord(t *testing.T) {
	t.Run("empty-correct", func(t *testing.T) {
		tests := []struct {
			word string
			want bool
		}{
			{"batty", true},
			{"fatty", false},
			{"party", false},
		}
		for i, test := range tests {
			h := History{
				almostLetters: []rune{'t', 'a', 't'},
				prohibitedLetters: [numLetters]char_set.CharSet{
					0: newCharSetHelper(t, 'f'),
				},
			}
			if want, got := test.want, h.allows(test.word); want != got {
				t.Errorf("test %v (shared history): wanted %v, got %v", i, want, got)
			}
		}
	})
	t.Run("custom-histories", func(t *testing.T) {
		tests := []struct {
			History
			name string
			word string
			want bool
		}{
			{
				name: `result{guess:"treat",score:"nannc"}`,
				History: History{
					correctLetters: [numLetters]rune{
						4: 't',
					},
					almostLetters: []rune{'r', 't'},
					prohibitedLetters: [numLetters]char_set.CharSet{
						newCharSetHelper(t, 't', 'e', 'a'),
						newCharSetHelper(t, 't', 'e', 'a', 'r'),
						newCharSetHelper(t, 't', 'e', 'a'),
						newCharSetHelper(t, 't', 'e', 'a'),
						newCharSetHelper(t, 't', 'e', 'a'),
					},
				},
				word: "robot",
				want: true,
			},
			{
				name: `result{guess:"shove",score: "accnc"}`,
				History: History{
					correctLetters: [numLetters]rune{
						1: 'h',
						2: 'o',
						4: 'e',
					},
					almostLetters: []rune{'s', 'h', 'o', 'e'},
					prohibitedLetters: [numLetters]char_set.CharSet{
						newCharSetHelper(t, 's', 'v'),
						newCharSetHelper(t, 'v'),
						newCharSetHelper(t, 'v'),
						newCharSetHelper(t, 'v'),
						newCharSetHelper(t, 'v'),
					},
				},
				word: "holes",
				want: false,
			},
		}
		for i, test := range tests {
			if want, got := test.want, test.History.allows(test.word); want != got {
				t.Errorf("test %v (%v) (with custom history): wanted %v, got %v", i, test.name, want, got)
			}
		}

	})
}

func newCharSetHelper(t *testing.T, chars ...rune) char_set.CharSet {
	t.Helper()
	var cs char_set.CharSet
	for _, ch := range chars {
		cs.Add(ch)
	}
	return cs
}

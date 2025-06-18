package main

import (
	"reflect"
	"testing"
)

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
		correctLetters: [numLetters]rune{
			4: 'y',
		},
		almostLetters: []rune{'a', 'y'},
		prohibitedLetters: [numLetters]charSet{
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't', 'a'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't'),
			newCharSetHelper(t, 'n', 's', 't'),
		},
	}
	wantWords := words{
		"alley": {},
		"ready": {},
	}
	var h history
	got := h
	got.addResult(r, &allWords)
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
		history
		guess
		score
		want history
	}{
		{
			guess: "treat",
			score: "nannc",
			want: history{
				correctLetters: [numLetters]rune{
					4: 't',
				},
				almostLetters: []rune{'r', 't'},
				prohibitedLetters: [numLetters]charSet{
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a', 'r'),
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a'),
					newCharSetHelper(t, 't', 'e', 'a'),
				},
			},
		},
		{
			guess: "shove",
			score: "accnc",
			want: history{
				correctLetters: [numLetters]rune{
					1: 'h',
					2: 'o',
					4: 'e',
				},
				almostLetters: []rune{'s', 'h', 'o', 'e'},
				prohibitedLetters: [numLetters]charSet{
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
		r := result{
			guess: test.guess,
			score: test.score,
		}
		want := test.want
		got := test.history
		got.mergeResult(r)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("test %v (for guess %q): histories not equal:\nwanted: %+v\ngot:    %+v", i, test.guess, want, got)
		}
	}
}

func TestString(t *testing.T) {
	h := history{
		correctLetters: [numLetters]rune{
			4: 'q',
		},
		almostLetters: []rune{'c', 'a', 'b'},
		prohibitedLetters: [numLetters]charSet{
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
		history
		newRequiredLetters []rune
		want               history
	}{
		{
			want: history{almostLetters: []rune{}},
		},
		{
			newRequiredLetters: []rune{'a', 'b', 'c'},
			want:               history{almostLetters: []rune{'a', 'b', 'c'}},
		},
		{
			history:            history{almostLetters: []rune{'a', 'b'}},
			newRequiredLetters: []rune{'a', 'a'},
			want:               history{almostLetters: []rune{'a', 'b', 'a'}},
		},
		{
			history:            history{almostLetters: []rune{'a', 'a', 'a'}},
			newRequiredLetters: []rune{'a', 'a', 'b', 'b', 'c'},
			want:               history{almostLetters: []rune{'a', 'a', 'a', 'b', 'b', 'c'}}, // this will prohibit all words
		},
	}
	for i, test := range tests {
		got := test.history
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
			h := history{
				almostLetters: []rune{'t', 'a', 't'},
				prohibitedLetters: [numLetters]charSet{
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
			history
			name string
			word string
			want bool
		}{
			{
				name: `result{guess:"treat",score:"nannc"}`,
				history: history{
					correctLetters: [numLetters]rune{
						4: 't',
					},
					almostLetters: []rune{'r', 't'},
					prohibitedLetters: [numLetters]charSet{
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
				history: history{
					correctLetters: [numLetters]rune{
						1: 'h',
						2: 'o',
						4: 'e',
					},
					almostLetters: []rune{'s', 'h', 'o', 'e'},
					prohibitedLetters: [numLetters]charSet{
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
			if want, got := test.want, test.history.allows(test.word); want != got {
				t.Errorf("test %v (%v) (with custom history): wanted %v, got %v", i, test.name, want, got)
			}
		}

	})
}

package server

import (
	"slices"
	"testing"

	"github.com/jacobpatterson1549/wordle-cheater/internal/letter_boxed"
)

const (
	letterBoxedParam = "letters"
)

func TestNewLetterBoxedCheater(t *testing.T) {
	tests := []struct {
		name      string
		query     map[string][]string
		wordsText string
		wantErr   bool
		want      LetterBoxedCheater
	}{
		{},
		{
			name:    "uneven letter count",
			query:   map[string][]string{letterBoxedParam: {"oddity"}},
			wantErr: true,
		},
		{
			name:    "multiple letters params",
			query:   map[string][]string{letterBoxedParam: {"eokmpjuarlcb", "eokmpjuarlcb"}},
			wantErr: true,
		},
		{
			name:      "ok",
			query:     map[string][]string{letterBoxedParam: {"eokmpjuarlcb"}},
			wordsText: "bore jock queen",
			want: LetterBoxedCheater{
				LetterBox: letter_boxed.LetterBox{
					Letters: "eokmpjuarlcb",
				},
				Result: letter_boxed.Result{
					Words: []string{"bore", "jock"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewLetterBoxedCheater(test.query, test.wordsText)
			switch {
			case got == nil, err != nil:
				if !test.wantErr {
					t.Errorf("unwanted error: %v", err)
				}
			case test.wantErr:
				t.Error("wanted error")
			case test.want.Letters != got.Letters:
				t.Errorf("letters not equal: \n wanted: %v \n    got: %v", test.want.Letters, got.Letters)
			case !slices.Equal(test.want.Words, got.Words):
				t.Errorf("words not equal: \n wanted: %v \n    got: %v", test.want.Words, got.Words)
			}
		})
	}
}

func TestLetterBoxedCheaterSortWords(t *testing.T) {
	words := []string{
		"apple",
		"cherry",
		"banana",
		"strawberry",
		"grape",
	}
	want := []string{
		"strawberry",
		"banana",
		"cherry",
		"apple",
		"grape",
	}
	var lbc LetterBoxedCheater
	slices.SortFunc(words, lbc.sortWords)
	got := words
	if !slices.Equal(want, got) {
		t.Errorf("not equal: \n wanted: %v \n    got: %v", want, got)
	}
}

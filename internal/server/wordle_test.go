package server

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/result"
)

func TestRunWordleCheater(t *testing.T) {
	tests := []struct {
		name   string
		query  map[string][]string
		wantOk bool
		want   WordleCheater
	}{
		{
			name:   "empty",
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{},
				},
			},
		},
		{
			name: "unexpected param",
			query: map[string][]string{
				"unknown": {"ignore"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{{}},
			},
		},
		{
			name: "invalid guess",
			query: map[string][]string{
				"g0":           {"word"},
				"s0":           {"ccccn"},
				"ShowPossible": {""},
			},
		},
		{
			name: "invalid score",
			query: map[string][]string{
				"g0":           {"words"},
				"s0":           {"right"},
				"ShowPossible": {""},
			},
		},
		{
			name: "one guess",
			query: map[string][]string{
				"g0":           {"forts"},
				"s0":           {"ccccn"},
				"g1":           {""},
				"s1":           {""},
				"ShowPossible": {""},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccn"},
					{},
				},
				Possible:     []string{"forte", "forth", "forty"},
				ShowPossible: true,
			},
		},
		{
			name: "two guesses",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccn"},
				"g1": {"forth"},
				"s1": {"ccccn"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccn"},
					{Guess: "forth", Score: "ccccn"},
					{},
				},
			},
		},
		{
			name: "first guess correct",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccc"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccc"},
				},
				Done: true,
			},
		},
		{
			name: "result fields not zero-indexed",
			query: map[string][]string{
				"g1": {"forts"},
				"s1": {"ccccn"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccn"},
					{},
				},
			},
		},
		{
			name: "extra guess",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccn"},
				"g1": {"forts"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccn"},
					{},
				},
			},
		},
		{
			name: "extra score",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccn"},
				"s1": {"ccccn"},
			},
			wantOk: true,
			want: WordleCheater{
				Results: []result.Result{
					{Guess: "forts", Score: "ccccn"},
					{},
				},
			},
		},
		{
			name: "duplicate guesses",
			query: map[string][]string{
				"g0": {"forts", "forth"},
				"s0": {"ccccn"},
			},
		},
		{
			name: "duplicate results",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccn", "ccccn"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			words := "forte forth forts forty"
			got, err := NewWordleCheater(test.query, words)
			switch {
			case !test.wantOk:
				if err == nil {
					t.Error("wanted error")
				}
			case err != nil:
				t.Errorf("unwanted error: %v", err)
			case got == nil, !reflect.DeepEqual(test.want, *got):
				t.Errorf("unequal: \n wanted: %+v \n got:    %+v", test.want, *got)
			}
		})
	}
}

func TestRunWordleCheaterBadWordsText(t *testing.T) {
	if _, err := NewWordleCheater(map[string][]string{}, "Words"); err == nil {
		t.Errorf("wanted error running with capitalized word")
	}
}

func TestRunWordleCheaterGuessCount(t *testing.T) {
	t.Run("guessCount", func(t *testing.T) {
		tests := []struct {
			n        int
			wantDone bool
		}{
			{0, false},
			{6, false},
			{9, true},
			{10, true},
		}
		for _, test := range tests {
			t.Run(strconv.Itoa(test.n), func(t *testing.T) {
				query := make(map[string][]string, test.n)
				for i := 0; i < test.n; i++ {
					query["g"+strconv.Itoa(i)] = []string{"xxxx" + string(byte('a'+i))}
					query["s"+strconv.Itoa(i)] = []string{"nnnnn"}
				}
				wordsText := "xxxxa xxxxb xxxxc xxxxd xxxxe xxxxf xxxxg xxxxh xxxxi xxxxj"
				got, err := NewWordleCheater(query, wordsText)
				switch {
				case err != nil:
					t.Errorf("unwanted error: %v", err)
				case got == nil, test.wantDone != got.Done:
					t.Errorf("done states: wanted %v, got %v", test.wantDone, got.Done)
				}
			})
		}
	})
}

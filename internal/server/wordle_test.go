package server

import (
	"reflect"
	"strconv"
	"testing"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/result"
)

func TestRunWordleCheaterPublic(t *testing.T) {
	o := words.WordsTextFile
	defer func() { words.WordsTextFile = o }()
	tests := []struct {
		words  string
		wantOk bool
	}{
		{wantOk: true},
		{words: "words", wantOk: true},
		{words: "Words"},
	}
	for _, test := range tests {
		t.Run(test.words, func(t *testing.T) {
			words.WordsTextFile = test.words
			got, err := RunWordleCheater(nil)
			switch {
			case err != nil:
				if test.wantOk {
					t.Error(err)
				}
			case !test.wantOk:
				t.Errorf("wanted error")
			case got == nil:
				t.Errorf("wanted cheater: %v", got)
			}
		})
	}
}

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
				"unknown": nil,
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
		},
		{
			name: "extra score",
			query: map[string][]string{
				"g0": {"forts"},
				"s0": {"ccccn"},
				"s1": {"ccccn"},
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
			got, err := runWordleCheater(test.query, words)
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

func TestRunWordleCheaterGuesssCount(t *testing.T) {
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
				wordsTextFile := "xxxxa xxxxb xxxxc xxxxd xxxxe xxxxf xxxxg xxxxh xxxxi xxxxj"
				got, err := runWordleCheater(query, wordsTextFile)
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

package server

import (
	"reflect"
	"testing"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/spelling_bee"
)

func TestRunSpellingBeeCheater(t *testing.T) {
	tests := []struct {
		name   string
		query  map[string][]string
		wantOk bool
		want   SpellingBeeCheater
	}{
		{
			name:   "empty",
			wantOk: true,
			want: SpellingBeeCheater{
				SpellingBee: spelling_bee.SpellingBee{
					MinLength: 4,
				},
			},
		},
		{
			name: "ok",
			query: map[string][]string{
				centralLetterParam: {"a"},
				otherLettersParam:  {"bcdefg"},
			},
			wantOk: true,
			want: SpellingBeeCheater{
				SpellingBee: spelling_bee.SpellingBee{
					CentralLetter: 'a',
					OtherLetters:  "bcdefg",
					MinLength:     4,
				},
			},
		},
		{
			name: "missing central-letter",
			query: map[string][]string{
				otherLettersParam: {"bcdefg"},
			},
		},
		{
			name: "missing other-letters",
			query: map[string][]string{
				centralLetterParam: {"a"},
			},
		},
		{
			name: "long central-letter",
			query: map[string][]string{
				centralLetterParam: {"aa"},
				"otherLetters":     {"bcdefg"},
			},
		},
		{
			name: "extra central-letter",
			query: map[string][]string{
				centralLetterParam: {"a", "a"},
				"otherLetters":     {"bcdefg"},
			},
		},
		{
			name: "long other-letters",
			query: map[string][]string{
				centralLetterParam: {"a"},
				otherLettersParam:  {"bcdefgh"},
			},
		},
		{
			name: "extra other-letters",
			query: map[string][]string{
				centralLetterParam: {"a"},
				otherLettersParam:  {"bcdefg", ""},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := RunSpellingBeeCheater(test.query)
			switch {
			case err != nil:
				if test.wantOk {
					t.Errorf("unwanted error: %v", err)
				}
			case !test.wantOk:
				t.Errorf("wanted error")
			case !reflect.DeepEqual(test.want, *got):
				t.Errorf("not equal: \n wanted: %v \n    got: %v", test.want, *got)
			}
		})
	}
}

func TestSpellingBeeCheaterSummary(t *testing.T) {
	o := words.WordsTextFile
	defer func() { words.WordsTextFile = o }()
	words.WordsTextFile = "bad apple yam may hi an my am a mamy"
	sb := spelling_bee.SpellingBee{
		CentralLetter: 'a',
		OtherLetters:  "my",
		MinLength:     2,
	}
	c := SpellingBeeCheater{sb}
	details := "PANGRAM!"
	want := Summary{
		TotalScore:   20,
		PangramCount: 3,
		Words: []Word{
			{Score: 7, Value: "mamy", Details: details},
			{Score: 6, Value: "yam", Details: details},
			{Score: 6, Value: "may", Details: details},
			{Score: 1, Value: "am"},
		},
	}
	got := c.Summary()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("not equal: \n wanted: %v \n    got: %v", want, got)
	}
}

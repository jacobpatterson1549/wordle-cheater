package spelling_bee

import (
	"slices"
	"testing"
)

func TestGetScores(t *testing.T) {
	tests := []struct {
		name          string
		sb            SpellingBee
		wordsTextFile string
		want          []Word
	}{
		{},
		{
			name:          "20250621",
			sb:            SpellingBee{CentralLetter: 'e', OtherLetters: "hcking", MinLength: 4},
			wordsTextFile: "stuff inching chicken checking hen nice electro",
			want: []Word{
				{Score: 1, Value: "nice"},
				{Score: 7, Value: "chicken"},
				{Score: 15, Value: "checking", IsPangram: true},
			},
		},
		{
			name:          "trimOtherLetters",
			sb:            SpellingBee{CentralLetter: 'f', OtherLetters: "nun"},
			wordsTextFile: "fun",
			want: []Word{
				{Score: 6, Value: "fun", IsPangram: true},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.sb.Words(test.wordsTextFile)
			if want, got := test.want, got; !slices.Equal(want, got) {
				t.Errorf("not equal: \n wanted: %v \n    got: %v", want, got)
			}
		})
	}
}

func TestWordLess(t *testing.T) {
	tests := []struct {
		name string
		a    Word
		b    Word
		want int
	}{
		{},
		{"byScore", Word{Score: 1, Value: "large"}, Word{Score: 8, Value: "large"}, -1},
		{"byPangramA", Word{IsPangram: true}, Word{}, 1},
		{"byPangramB", Word{}, Word{IsPangram: true}, -1},
		{"byLength", Word{Score: 1, Value: "large"}, Word{Score: 1, Value: "tiny"}, 1},
		{"byValue", Word{Score: 1, Value: "long"}, Word{Score: 1, Value: "acne"}, 1},
		{"equal", Word{Score: 1, Value: "long"}, Word{Score: 1, Value: "long"}, 0},
	}
	sign := func(i int) int {
		if i < 0 {
			return -1
		}
		if i > 0 {
			return 1
		}
		return 0
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := wordLess(test.a, test.b)
			got = sign(got)
			if want, got := test.want, got; want != got {
				t.Errorf("wanted %v, got %v", want, got)
			}
		})
	}
}

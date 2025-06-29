package letter_boxed

import (
	"slices"
	"testing"
)

func TestWords(t *testing.T) {
	tests := []struct {
		name      string
		wordsText string
		lb        LetterBox
		wantOk    bool
		want      []string
	}{
		{},
		{
			name: "no letters",
			lb:   LetterBox{Letters: "", BoxSideCount: 3, MinWordLength: 1},
		},
		{
			name: "bad box side count",
			lb:   LetterBox{Letters: "abc", MinWordLength: 1},
		},
		{
			name: "bad min word length",
			lb:   LetterBox{Letters: "ham", BoxSideCount: 3, MinWordLength: -1},
		},
		{
			name: "uneven letter count",
			lb:   LetterBox{Letters: "rats", BoxSideCount: 3, MinWordLength: 1},
		},
		{
			name:      "simple",
			wordsText: "ab cab bad",
			lb:        LetterBox{Letters: "abc", BoxSideCount: 3, MinWordLength: 2},
			wantOk:    true,
			want:      []string{"ab", "cab"},
		},
		{
			name:      "20250622",
			wordsText: "zebra but fickle eat tamp puck bike left limp",
			lb:        LetterBox{Letters: "lmikfaecputb", BoxSideCount: 4, MinWordLength: 3},
			wantOk:    true,
			want:      []string{"bike", "eat", "fickle", "left", "puck", "tamp"},
		},
		{
			name:      "20250622",
			wordsText: "zebra but fickle eat tamp puck bike left limp",
			lb:        LetterBox{Letters: "lmikfaecputb", BoxSideCount: 4, MinWordLength: 3},
			wantOk:    true,
			want:      []string{"bike", "eat", "fickle", "left", "puck", "tamp"},
		},
		{
			name:      "two letters",
			wordsText: "odd dodo",
			lb:        LetterBox{Letters: "do", BoxSideCount: 2, MinWordLength: 3},
			wantOk:    true,
			want:      []string{"dodo"},
		},
		{
			name:      "duplicate letters",
			wordsText: "a aa aaa",
			lb:        LetterBox{Letters: "aaaa", BoxSideCount: 4, MinWordLength: 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.lb.Words(test.wordsText)
			if err != nil {
				if test.want != nil {
					t.Errorf("unwanted error: %v", err)
				}
				return
			}
			if want, got := test.want, got; !slices.Equal(want, got) {
				t.Errorf("not equal: \n wanted: %v \n    got: %v", want, got)
			}
		})
	}
}

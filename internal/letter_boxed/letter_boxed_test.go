package letter_boxed

import (
	"slices"
	"testing"
)

func TestWords(t *testing.T) {
	tests := []struct {
		name         string
		wordsText    string
		letterGroups []string
		want         []string
	}{
		{
			name:         "20250622",
			wordsText:    "zebra but fickle eat tamp puck bike left limp",
			letterGroups: []string{"lmi", "kfa", "ecp", "utb"},
			want:         []string{"bike", "eat", "fickle", "left", "puck", "tamp"},
		},
		{
			name:         "two letters",
			wordsText:    "odd dodo",
			letterGroups: []string{"d", "o"},
			want:         []string{"dodo"},
		},
		{
			name:         "duplicate letters",
			wordsText:    "a aa aaa",
			letterGroups: []string{"aa"},
		},
		{
			name:         "duplicate letters across groups",
			wordsText:    "a aa aaa",
			letterGroups: []string{"a", "a"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Words(test.wordsText, test.letterGroups)
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

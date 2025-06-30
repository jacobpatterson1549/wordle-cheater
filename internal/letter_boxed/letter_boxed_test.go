package letter_boxed

import (
	"slices"
	"testing"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/char_set"
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
			got, err := test.lb.words(test.wordsText)
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

func TestNewConnection(t *testing.T) {
	var c connections
	for _, r := range "abcde" {
		c.targets.Add(r)
	}
	want := connection{
		Word: "abracadabra",
	}
	want.targets.Add('a')
	want.targets.Add('b')
	want.targets.Add('c')
	want.targets.Add('d')
	got := c.newConnection(want.Word)
	if want != got {
		t.Errorf("not equal: \n wanted: %v \n    got: %v", want, got)
	}
}

func TestConnections(t *testing.T) {
	tests := []struct {
		name      string
		lb        LetterBox
		wantErr   bool
		wantWords []string
		wantNumC  int
	}{
		{},
		{
			name: "bad config",
			lb: LetterBox{
				Letters: "x",
			},
			wantErr: true,
		},
		{	
			name: "one word",
			lb: LetterBox{
				Letters:       "love",
				BoxSideCount:  4,
				MinWordLength: 4,
			},
			wantWords: []string{"love"},
			wantNumC:  1,
		},
		{
			name: "no connections",
			lb: LetterBox{
				Letters:      "farwide",
				BoxSideCount: 7,
				MinWordLength: 3,
			},
			wantWords: []string{"far", "wide"},
			wantNumC:  2,
		},
	}
	wordsText := "far wide eat sleep pray love vowel"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.lb.Connections(wordsText)
			switch {
			case err != nil:
				if !test.wantErr {
					t.Errorf("unwanted error: %v", err)
				}
			case test.wantErr, got == nil:
				t.Error("wanted error")
			case !slices.Equal(test.wantWords, got.Words):
				t.Errorf("words not equal: \n wanted: %v \n    got: %v", test.wantWords, got.Words)
			case test.wantNumC != len(got.Connections):
				t.Errorf("connections counts not equal: \n wanted: %v \n    got: %v", test.wantNumC, len(got.Connections))
			}
		})
	}
}

func TestNewConnections(t *testing.T) {
	wordsText := "apple pear plum grape"
	t.Run("ok", func(t *testing.T) {
		lb := LetterBox{
			Letters:       "abcdefghijklmnopqrstuvwxy",
			BoxSideCount:  25,
			MinWordLength: 4,
		}
		got, err := lb.newConnections(wordsText)
		switch {
		case err != nil, got == nil:
			t.Errorf("unwanted error: %v", err)
		case got.targets.Length() != 25:
			t.Errorf("target counts: %v", got.targets.Length())
		case len(got.words) != 3:
			t.Errorf("word count: %v", len(got.all))
		case len(got.all) != 3: // excludes repeated letters
			t.Errorf("connection count: %v", len(got.all))
		case len(got.startsWith['p'-'a']) != 2:
			t.Errorf("connection starts with 'p' count: %v", len(got.startsWith['p'-'a']))
		case len(got.endsWith['e'-'a']) != 1:
			t.Errorf("connection ends with 'e' count: %v", len(got.endsWith['e'-'a']))
		}
	})
	t.Run("bad config", func(t *testing.T) {
		lb := LetterBox{
			Letters: "x",
		}
		_, err := lb.newConnections(wordsText)
		if err == nil {
			t.Errorf("wanted error with invalid letterBox")
		}
	})
}

func TestConnectionLess(t *testing.T) {
	var target1, target2 char_set.CharSet
	target1.Add('a')
	target2.Add('a')
	target2.Add('b')
	tests := []struct {
		name string
		a    connection
		b    connection
		want int
	}{
		{
			name: "by target count",
			a:    connection{Word: "zzz", targets: target2},
			b:    connection{Word: "aaa", targets: target1},
			want: -1,
		},
		{
			name: "by word length",
			a:    connection{Word: "aaba", targets: target1},
			b:    connection{Word: "zab", targets: target1},
			want: +1,
		},
		{
			name: "by alphabetical",
			a:    connection{Word: "aaa", targets: target1},
			b:    connection{Word: "aab", targets: target1},
			want: -1,
		},
		{
			name: "equal",
			a:    connection{Word: "aab", targets: target1},
			b:    connection{Word: "aab", targets: target1},
			want: 0,
		},
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
			got := connectionLess(test.a, test.b)
			if want, got := test.want, sign(got); want != got {
				t.Errorf("wanted %v, got %v", want, got)
			}
			got = connectionLess(test.b, test.a)
			if want, got := -test.want, sign(got); want != got {
				t.Errorf("flipped: wanted %v, got %v", want, got)
			}
		})
	}
}

func TestRemoveConnection(t *testing.T) {
	var targetsA, targetsAB, targetsB char_set.CharSet
	targetsA.Add('a')
	targetsAB.Add('a')
	targetsAB.Add('b')
	targetsB.Add('b')
	cwHat := connection{
		Word:    "hat",
		targets: targetsA,
	}
	cwBat := connection{
		Word:    "bat",
		targets: targetsAB,
	}
	c := connections{
		targets: targetsAB,
		words:   []string{"bat", "hat"},
		all:     []connection{cwHat, cwBat},
		startsWith: [26][]*connection{
			'b' - 'a': {&cwBat},
			'h' - 'a': {&cwHat},
		},
		endsWith: [26][]*connection{
			't' - 'a': {&cwHat, &cwBat},
		},
	}
	c.remove(cwHat)
	switch {
	case c.targets != targetsB:
		t.Errorf("all targets: wanted %v, got %v", targetsB, c.targets)
	case len(c.words) != 2:
		t.Errorf("both words should still exist: %v", c.words)
	case len(c.all) != 2:
		t.Errorf("both connections should still exist: %v", c.all)
	case len(c.startsWith['b'-'a']) != 1:
		t.Errorf("starts with length: %v", len(c.startsWith['b'-'a']))
	default:
		for _, cw := range c.all {
			if cw.targets.Has('h') {
				t.Errorf("target had 'h': %v", cw.targets)
			}
		}
	}
}

func TestMinSubSet(t *testing.T) {
	// eokmpjuarlcb => corporeal lumberjack lumber jack
	// nvaoguihltrd => individuating grot that told grunt gator
}

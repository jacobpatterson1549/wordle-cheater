package letter_boxed

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/char_set"
)

type (
	LetterBox struct {
		Letters       string
		BoxSideCount  int
		MinWordLength int
	}
	Result struct {
		Words       []string
		Connections []string
	}
	groups     map[rune]int
	connection struct {
		Word    string
		targets char_set.CharSet
	}
	connectionHeap []*connection
	solver         struct {
		targets     char_set.CharSet
		words       []string
		all         []connection
		startsWith  [26][]*connection
		endsWith    [26][]*connection
		targetFreqs [26]int
	}
)

func (lb LetterBox) words(wordsText string) ([]string, error) {
	letters := []rune(lb.Letters)
	switch {
	case len(letters) == 0:
		return nil, nil
	case lb.BoxSideCount <= 0:
		return nil, fmt.Errorf("wanted positive box side count: %v", lb.BoxSideCount)
	case lb.MinWordLength <= 0:
		return nil, fmt.Errorf("wanted positive required word length: %v", lb.MinWordLength)
	case len(letters)%lb.BoxSideCount != 0:
		return nil, fmt.Errorf("letters on each side of box not equal")
	}
	words := strings.Fields(wordsText)
	letterGroups := make([]string, lb.BoxSideCount)
	k := len(letters) / lb.BoxSideCount
	for i := range lb.BoxSideCount {
		j := i * k
		letterGroups[i] = string(letters[j : j+k])
	}
	g, err := newGroups(letterGroups)
	if err != nil {
		return nil, err
	}
	var validWords []string
	for _, word := range words {
		if len(word) >= lb.MinWordLength && g.allows(word) {
			validWords = append(validWords, word)
		}
	}
	slices.Sort(validWords)
	return validWords, nil
}

func newGroups(letterGroups []string) (*groups, error) {
	g := make(groups)
	for key, side := range letterGroups {
		for _, r := range side {
			if _, ok := g[r]; ok {
				return nil, fmt.Errorf("%q in duplicated or in multiple groups", string(r))
			}
			g[r] = key
		}
	}
	return &g, nil
}

func (g groups) allows(word string) bool {
	var prevKey int
	for i, r := range word {
		currKey, ok := g[r]
		switch {
		case !ok,
			i != 0 && prevKey == currKey:
			return false
		}
		prevKey = currKey
	}
	return len(word) > 0
}

func (lb LetterBox) Solve(wordsText string) (*Result, error) {
	// TODO: write solver, add timeout if algorithm is slow
	words, err := lb.words(wordsText)
	if err != nil {
		return nil, err
	}
	r := Result{
		Words: words,
	}
	return &r, nil
}

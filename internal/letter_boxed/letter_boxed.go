package letter_boxed

import (
	"fmt"
	"slices"
	"strings"
)

type groups map[rune]int

func Words(wordsText string, letterGroups []string) ([]string, error) {
	words := strings.Fields(wordsText)
	g, err := newGroups(letterGroups)
	if err != nil {
		return nil, err
	}
	var validWords []string
	for _, word := range words {
		if g.allows(word) {
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

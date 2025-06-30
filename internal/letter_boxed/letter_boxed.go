package letter_boxed

import (
	"cmp"
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
	connections struct {
		targets    char_set.CharSet
		words      []string
		all        []connection
		startsWith [26][]*connection
		endsWith   [26][]*connection
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

func (lb LetterBox) Connections(wordsText string) (*Result, error) {
	c, err := lb.newConnections(wordsText)
	if err != nil {
		return nil, fmt.Errorf("finding connections: %w", err)
	}
	minConnections, err := c.minSubSet()
	if err != nil {
		minConnections = append(minConnections, fmt.Sprintf("[%v]", err))
	}
	r := Result{
		Words:       c.words,
		Connections: minConnections,
	}
	return &r, nil
}

func (lb LetterBox) newConnections(wordsText string) (*connections, error) {
	words, err := lb.words(wordsText)
	if err != nil {
		return nil, err
	}
	var c connections
	c.initTargets(lb)
	c.words = words
	c.initAll()
	c.initStartsEndsWith()
	return &c, nil
}

func (c *connections) initTargets(lb LetterBox) {
	for _, r := range lb.Letters {
		c.targets.Add(r)
	}
}

func (c *connections) initAll() {
	c.all = make([]connection, len(c.words))
	for i, w := range c.words {
		c.all[i] = c.newConnection(w)
	}
}

func (c *connections) initStartsEndsWith() {
	for i := range c.all {
		letters := []rune(c.all[i].Word)
		j := letters[0] - 'a'
		k := letters[len(letters)-1] - 'a'
		c.startsWith[j] = append(c.startsWith[j], &c.all[i])
		c.endsWith[k] = append(c.endsWith[k], &c.all[i])
	}
}

func (c connections) newConnection(word string) connection {
	cw := connection{
		Word: word,
	}
	for _, r := range word {
		if c.targets.Has(r) {
			cw.targets.Add(r)
		}
	}
	return cw
}

func (c connections) minSubSet() ([]string, error) {
	missCount, maxMissCount := 0, 5
	var mss []string
	for len(c.all) > 0 && c.targets.Length() > 0 {
		targetsLength := c.targets.Length()
		if missCount >= maxMissCount {
			return mss, fmt.Errorf("could not finish building connections")
		}
		mss = c.removeMin(mss)
		if c.targets.Length() == targetsLength {
			missCount++
		}
		// fmt.Println("CONNECTIONS:", mss, c.targets.String())
	}
	return mss, nil
}

func (c *connections) removeMin(mss []string) []string {
	switch {
	case len(mss) == 0 && len(c.all) > 0:
		m := slices.MinFunc(c.all, connectionLess)
		c.remove(m)
		mss = append(mss, m.Word)
	default:
		start, end := mss[0], mss[len(mss)-1]
		cwStart, okStart := connectionMinRef(c.endsWith[start[0]-'a'])
		cwEnd, okEnd := connectionMinRef(c.startsWith[end[len(end)-1]-'a'])
		switch {
		case okStart && (!okEnd || connectionLessRef(cwStart, cwEnd) < 0):
			c.remove(*cwStart)
			mss = append(mss, "")
			copy(mss[1:], mss)
			mss[0] = cwStart.Word
		case okEnd:
			c.remove(*cwEnd)
			mss = append(mss, cwEnd.Word)
		}
	}
	return mss
}

func (c *connections) remove(cw connection) {
	// fmt.Println("removing ", cw.Word, cw.targets.String())
	for _, r := range cw.Word {
		c.targets.Remove(r)
		for _, cw := range c.all {
			cw.targets.Remove(r)
		}
	}
}

func connectionLess(a, b connection) int {
	if a, b := a.targets.Length(), b.targets.Length(); a != b {
		return b - a
	}
	if a, b := len(a.Word), len(b.Word); a != b {
		return a - b
	}
	return cmp.Compare(a.Word, b.Word)
}

func connectionMinRef(s []*connection) (cw *connection, ok bool) {
	if len(s) == 0 {
		return nil, false
	}
	return slices.MinFunc(s, connectionLessRef), true
}

func connectionLessRef(a, b *connection) int {
	return connectionLess(*a, *b)
}

package char_set

import (
	"fmt"
	"strings"
)

// CharSet is a bit field that stores the letters a-z
type CharSet uint32

// Add includes the character to the set.  Panics if the character is not in a-z
func (cs *CharSet) Add(ch rune) {
	if !cs.valid(ch) {
		panic(fmt.Errorf("%c is not in a-z", ch))
	}
	*cs |= cs.singleton(ch)
}

func (cs *CharSet) AddAll(s string) {
	for _, r := range s {
		cs.Add(r)
	}
}

func (cs *CharSet) Remove(ch rune) {
	if !cs.valid(ch) {
		panic(fmt.Errorf("%c is not in a-z", ch))
	}
	*cs &^= cs.singleton(ch)
}

func (cs *CharSet) RemoveAll(s string) {
	for _, r := range s {
		cs.Remove(r)
	}
}

// Has determines if the character is in the set
func (cs CharSet) Has(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return (cs & cs.singleton(ch)) != 0
}

// AddWouldFill determines if the charset is filled with the letters a-z
func (cs CharSet) AddWouldFill(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return cs|cs.singleton(ch) == cs.singleton('z'+1)-1
}

// String creates a string of the characters in the set, in ascending order
func (cs CharSet) String() string {
	var b strings.Builder
	b.WriteRune('[')
	for ch := rune('a'); ch <= 'z'; ch++ {
		if cs.Has(ch) {
			b.WriteRune(ch)
		}
	}
	b.WriteRune(']')
	return b.String()
}

// valid determines if the byte can be used in the charSet, if it is a-z
func (CharSet) valid(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}

// singleton creates a singleton charSet from the character
func (CharSet) singleton(ch rune) CharSet {
	return 1 << (ch - 'a')
}

func (cs CharSet) Length() int {
	n := 0
	for ch := rune('a'); ch <= 'z'; ch++ {
		if cs.Has(ch) {
			n++
		}
	}
	return n
}

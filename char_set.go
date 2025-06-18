package main

import (
	"fmt"
	"strings"
)

// charSet is a bit field that stores the letters a-z
type charSet uint32

// Add includes the character to the set.  Panics if the character is not in a-z
func (cs *charSet) Add(ch rune) {
	if !cs.valid(ch) {
		panic(fmt.Errorf("%c is not in a-z", ch))
	}
	*cs |= cs.singleton(ch)
}

// Has determines if the character is in the set
func (cs charSet) Has(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return (cs & cs.singleton(ch)) != 0
}

// AddWouldFill determines if the charset is filled with the letters a-z
func (cs charSet) AddWouldFill(ch rune) bool {
	if !cs.valid(ch) {
		return false
	}
	return cs|cs.singleton(ch) == cs.singleton('z'+1)-1
}

// String creates a string of the characters in the set, in ascending order
func (cs charSet) String() string {
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
func (charSet) valid(ch rune) bool {
	return 'a' <= ch && ch <= 'z'
}

// singleton creates a singleton charSet from the character
func (charSet) singleton(ch rune) charSet {
	return 1 << (ch - 'a')
}

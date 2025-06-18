package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// words is a collection of unique strings
type words map[string]struct{}

// newWords loads the words from the file.
// Words are separated by whitespace (spaces/newlines).
// An error is returned if any words are not <<numLetters characters long and lowercase.
func newWords(a string) (*words, error) {
	lines := strings.Fields(a)
	m := make(words, len(lines))
	for _, w := range lines {
		if len(w) != numLetters {
			return nil, fmt.Errorf("wanted all words to be %v letters long, got %q", numLetters, w)
		}
		if w != strings.ToLower(w) {
			return nil, fmt.Errorf("wanted all words to be lowercase, got %q", w)
		}
		m[w] = struct{}{}
	}
	return &m, nil
}

// copy creates a new, identical duplication of the words
func (m words) copy() *words {
	m2 := make(words, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return &m2
}

// sorted combines and sorts the words into a csv string
func (m words) sorted() string {
	s := make([]string, 0, len(m))
	for w := range m {
		s = append(s, w)
	}
	sort.Strings(s)
	j := strings.Join(s, ",")
	return j
}

// scanShowPossible prompts to display the words
func (m words) scanShowPossible(rw io.ReadWriter) error {
	fmt.Fprintf(rw, "show possible words [Yn]: ")
	var choice string
	n, err := fmt.Fscanf(rw, "%s", &choice)
	if err == io.EOF || (n != 0 && err != nil) {
		return fmt.Errorf("scanning choice: %v", err)
	}
	choice = strings.ToLower(choice)
	if n != 0 && len(choice) > 0 && choice[0] == 'n' {
		return nil
	}
	fmt.Fprintf(rw, "remaining valid words: %v\n", m.sorted())
	return nil
}

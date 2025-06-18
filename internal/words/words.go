package words

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Words is a collection of unique strings
type Words map[string]struct{}

// New loads the words from the file.
// Words are separated by whitespace (spaces/newlines).
// An error is returned if any words are not <<numLetters characters long and lowercase.
func New(a string) (*Words, error) {
	lines := strings.Fields(a)
	m := make(Words, len(lines))
	for _, w := range lines {
		if len(w) != 5 {
			continue
		}
		if w != strings.ToLower(w) {
			return nil, fmt.Errorf("wanted all words to be lowercase, got %q", w)
		}
		m[w] = struct{}{}
	}
	return &m, nil
}

// Copy creates a new, identical duplication of the words
func (m Words) Copy() *Words {
	m2 := make(Words, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return &m2
}

// sorted combines and sorts the words into a csv string
func (m Words) sorted() string {
	s := make([]string, 0, len(m))
	for w := range m {
		s = append(s, w)
	}
	sort.Strings(s)
	j := strings.Join(s, ",")
	return j
}

// ScanShowPossible prompts to display the words
func (m Words) ScanShowPossible(rw io.ReadWriter) error {
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

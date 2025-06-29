package server

import "testing"

func TestUniquePageTitles(t *testing.T) {
	titles := []string{
		wordlePage.Title,
		spellingBeePage.Title,
		letterBoxedPage.Title,
	}
	m := make(map[string]struct{}, len(titles))
	for _, title := range titles {
		if _, ok := m[title]; ok {
			t.Errorf("duplicate page title, used as unique key: %v", title)
		}
		m[title] = struct{}{}
	}
	if want, got := len(titles), len(m); want != got {
		t.Errorf("distinct title count: wanted %v, got %v", want, got)
	}
}

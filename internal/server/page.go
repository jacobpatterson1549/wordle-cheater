package server

type (
	pageDisplay struct {
		page
		NoJS    bool
		Cheater any
	}
	page struct {
		Title      string
		tmplName   string
		newCheater func(query map[string][]string, wordsText string) (any, error)
	}
)

var (
	wordlePage = page{
		Title:      "Wordle Cheater",
		tmplName:   "wordle.html",
		newCheater: wrapCheater(NewWordleCheater),
	}
	spellingBeePage = page{
		Title:      "Spelling Bee Cheater",
		tmplName:   "spelling_bee.html",
		newCheater: wrapCheater(NewSpellingBeeCheater),
	}
	letterBoxedPage = page{
		Title:      "Letter Boxed Cheater",
		tmplName:   "letter_boxed.html",
		newCheater: wrapCheater(NewLetterBoxedCheater),
	}
)

func wrapCheater[T any](f func(query map[string][]string, wordsText string) (T, error)) func(query map[string][]string, wordsText string) (any, error) {
	return func(query map[string][]string, wordsText string) (any, error) {
		c, err := f(query, wordsText)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func (pt page) newPage(query map[string][]string, wordsText string) (*pageDisplay, error) {
	c, err := pt.newCheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	p := pageDisplay{
		page:    pt,
		Cheater: c,
	}
	return &p, nil
}

func (pt page) IsWordle() bool {
	return pt.Title == wordlePage.Title
}

func (pt page) IsSpellingBee() bool {
	return pt.Title == spellingBeePage.Title
}

func (pt page) IsLetterBoxed() bool {
	return pt.Title == letterBoxedPage.Title
}

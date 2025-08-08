package server

type (
	display struct {
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

func (p page) newDisplay(query map[string][]string, wordsText string) (*display, error) {
	c, err := p.newCheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	d := display{
		page:    p,
		Cheater: c,
	}
	return &d, nil
}

func (p page) IsWordle() bool {
	return p.Title == wordlePage.Title
}

func (p page) IsSpellingBee() bool {
	return p.Title == spellingBeePage.Title
}

func (p page) IsLetterBoxed() bool {
	return p.Title == letterBoxedPage.Title
}

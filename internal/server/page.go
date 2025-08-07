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
		Title:    "Wordle Cheater",
		tmplName: "wordle.html",
		newCheater: func(query map[string][]string, wordsText string) (any, error) {
			return NewWordleCheater(query, wordsText)
		},
	}
	spellingBeePage = page{
		Title:    "Spelling Bee Cheater",
		tmplName: "spelling_bee.html",
		newCheater: func(query map[string][]string, wordsText string) (any, error) {
			return NewSpellingBeeCheater(query, wordsText)
		},
	}
	letterBoxedPage = page{
		Title:    "Letter Boxed Cheater",
		tmplName: "letter_boxed.html",
		newCheater: func(query map[string][]string, wordsText string) (any, error) {
			return NewLetterBoxedCheater(query, wordsText)
		},
	}
)

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

package server

type (
	pageDisplay[C any] struct {
		page[C]
		NoJS    bool
		Cheater any
	}
	page[C any] struct {
		Title            string
		HtmxTemplateName string
		newCheater       func(query map[string][]string, wordsText string) (*C, error)
	}
)

var (
	wordlePage      = page[WordleCheater]{"Wordle Cheater", "wordle.html", NewWordleCheater}
	spellingBeePage = page[SpellingBeeCheater]{"Spelling Bee Cheater", "sbc-response", NewSpellingBeeCheater}
	letterBoxedPage = page[LetterBoxedCheater]{"Letter Boxed Cheater", "lbc-response", NewLetterBoxedCheater}
)

func (pt page[C]) newPage(query map[string][]string, wordsText string) (*pageDisplay[C], error) {
	c, err := pt.newCheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	p := pageDisplay[C]{
		page:    pt,
		Cheater: c,
	}
	return &p, nil
}

func (pt page[C]) IsWordle() bool {
	return pt.Title == wordlePage.Title
}

func (pt page[C]) IsSpellingBee() bool {
	return pt.Title == spellingBeePage.Title
}

func (pt page[C]) IsLetterBoxed() bool {
	return pt.Title == letterBoxedPage.Title
}

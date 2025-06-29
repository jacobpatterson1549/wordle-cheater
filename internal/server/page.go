package server

type (
	Page[C any] struct {
		pageType[C]
		NoJS    bool
		Cheater any
	}
	pageType[C any] struct {
		Title            string
		HtmxTemplateName string
		newCheater       func(query map[string][]string, wordsText string) (*C, error)
	}
)

var (
	wordleType      = pageType[WordleCheater]{"Wordle Cheater", "wordle.html", NewWordleCheater}
	spellingBeeType = pageType[SpellingBeeCheater]{"Spelling Bee Cheater", "sbc-response", NewSpellingBeeCheater}
	letterBoxedType = pageType[LetterBoxedCheater]{"Letter Boxed Cheater", "lbc-response", NewLetterBoxedCheater}
)

func (pt pageType[C]) newPage(query map[string][]string, wordsText string) (*Page[C], error) {
	c, err := pt.newCheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	p := Page[C]{
		pageType: pt,
		Cheater:  c,
	}
	return &p, nil
}

func (pt pageType[C]) IsWordle() bool {
	return pt.Title == wordleType.Title
}

func (pt pageType[C]) IsSpellingBee() bool {
	return pt.Title == spellingBeeType.Title
}

func (pt pageType[C]) IsLetterBoxed() bool {
	return pt.Title == letterBoxedType.Title
}

package server

type (
	Page struct {
		NoJS bool
		pageType
		Cheater any
	}
	pageType int
)

const (
	wordle_type pageType = iota
	spelling_bee_type
	letter_boxed_type
)

func (pt pageType) newPage(query map[string][]string, wordsText string) (*Page, error) {
	c, err := pt.newCheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	p := Page{
		pageType: pt,
		Cheater:  c,
	}
	return &p, nil
}

func (pt pageType) Title() string {
	switch pt {
	case spelling_bee_type:
		return "Spelling Bee Cheater"
	case letter_boxed_type:
		return "Letter Boxed"
	default:
		return "Wordle Cheater"
	}
}

func (pt pageType) newCheater(query map[string][]string, wordsText string) (any, error) {
	switch pt {
	case spelling_bee_type:
		return NewSpellingBeeCheater(query, wordsText)
	case letter_boxed_type:
		return NewLetterBoxedCheater(query, wordsText)
	default:
		return NewWordleCheater(query, wordsText)
	}
}

func (pt pageType) HtmxTemplateName() string {
	switch pt {
	case spelling_bee_type:
		return "sbc-response"
	case letter_boxed_type:
		return "lbc-response"
	default:
		return "wordle.html"
	}
}

func (pt pageType) IsWordle() bool {
	return pt == wordle_type
}

func (pt pageType) IsSpellingBee() bool {
	return pt == spelling_bee_type
}

func (pt pageType) IsLetterBoxed() bool {
	return pt == letter_boxed_type
}

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
)

func (pt pageType) newPage(query map[string][]string, wordsText string) (*Page, error) {
	cheater, err := pt.cheater(query, wordsText)
	if err != nil {
		return nil, err
	}
	p := Page{
		pageType: pt,
		Cheater:  cheater,
	}
	return &p, nil
}

func (pt pageType) Title() string {
	switch pt {
	case spelling_bee_type:
		return "Spelling Bee Cheater"
	default:
		return "Wordle Cheater"
	}
}

func (pt pageType) cheater(query map[string][]string, wordsText string) (any, error) {
	switch pt {
	case spelling_bee_type:
		return RunSpellingBeeCheater(query, wordsText)
	default:
		return RunWordleCheater(query, wordsText)
	}
}

func (pt pageType) HtmxTemplateName() string {
	switch pt {
	case spelling_bee_type:
		return "sbc-response"
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

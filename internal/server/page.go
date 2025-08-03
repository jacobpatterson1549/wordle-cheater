package server

type (
	pageDisplay[C any] struct {
		page[C]
		NoJS    bool
		Cheater any
	}
	page[C any] struct {
		Title        string
		tmplName     string
		newCheater   func(query map[string][]string, wordsText string) (*C, error)
		Instructions []string
	}
)

var (
	wordlePage = page[WordleCheater]{
		Title:      "Wordle Cheater",
		tmplName:   "wordle.html",
		newCheater: NewWordleCheater,
		Instructions: []string{
			"Wordle-Cheater is a word-guessing helper.",
			"Each guess must be five (5) letters long.",
			"Letters for each guess are assigned a score:",
			"- 'C' for correct - letter is in the word in the same position.",
			"- 'A' for almost - letter is in the word, but in a different position.",
			"- 'N' for not correct - letter is not in the word at all.",
			"Scores for guesses are cumulatively applied.",
			"Check the 'Show Possible' checkbox to see valid words after submitting another guess.",
		},
	}
	spellingBeePage = page[SpellingBeeCheater]{
		Title:      "Spelling Bee Cheater",
		tmplName:   "spelling_bee.html",
		newCheater: NewSpellingBeeCheater,
		Instructions: []string{
			"Spelling Bee Cheater finds words that match the pattern.",
			"The Central letter must be in each word.",
			"Words must be at least four (4) letters long.",
			"Other letters contains the list of six (6) distinct secondary letters that can be used.",
			"The score of a word is its letter count.",
			"However, short, four (4) letter words have a score of one (1).",
			"Words that use all the Other letters are Pangrams and get a bonus of seven (7) points.",
		},
	}
	letterBoxedPage = page[LetterBoxedCheater]{
		Title:      "Letter Boxed Cheater",
		tmplName:   "letter_boxed.html",
		newCheater: NewLetterBoxedCheater,
		Instructions: []string{
			"Letter Boxed Cheater lists words that can be formed by alternating letters in a box pattern.",
			"Each side of the box has three (3) letters.",
			"Words are be formed by jumping between box edges",
			"Enter letters for each side together, resulting in a distinct, twelve (12) letter box state.",
		},
	}
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

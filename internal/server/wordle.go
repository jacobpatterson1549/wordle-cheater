package server

import (
	"fmt"
	"slices"

	words "github.com/jacobpatterson1549/wordle-cheater"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/guess"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/result"
	"github.com/jacobpatterson1549/wordle-cheater/internal/wordle/score"
)

type WordleCheater struct {
	Results      []result.Result
	Possible     []string
	ShowPossible bool
	Done         bool
}

func RunWordleCheater(query map[string][]string, wordsText string) (*WordleCheater, error) {
	for k, v := range query {
		if len(v) != 1 {
			return nil, fmt.Errorf("wanted only one value for %q", k)
		}
	}

	m, err := words.New(wordsText)
	if err != nil {
		return nil, fmt.Errorf("creating word list: %w", err)
	}

	wc, err := newWordleCheater(query, *m)
	if err != nil {
		return nil, fmt.Errorf("parsing query: %w", err)
	}
	return wc, nil
}

func newWordleCheater(query map[string][]string, m words.Words) (*WordleCheater, error) {

	var wc WordleCheater
	var h result.History

	for i := range 10 {
		r, err := parseResult(query, i)
		switch {
		case err != nil:
			return nil, err
		case r != nil:
			h.AddResult(*r, &m)
			wc.Results = append(wc.Results, *r)
		}
	}

	if _, ok := query["ShowPossible"]; ok {
		wc.ShowPossible = true
	}

	wc.Done = len(wc.Results) >= 9 ||
		(len(wc.Results) > 0 && wc.Results[len(wc.Results)-1].Score == score.AllCorrect)
	if !wc.Done {
		wc.Results = append(wc.Results, result.Result{})
	}

	if wc.ShowPossible {
		wc.Possible = make([]string, 0, len(m))
		for k := range m {
			wc.Possible = append(wc.Possible, k)
		}
		slices.Sort(wc.Possible)
	}

	return &wc, nil
}

func parseResult(query map[string][]string, i int) (*result.Result, error) {
	guessKey := fmt.Sprintf("g%v", i)
	scoreKey := fmt.Sprintf("s%v", i)

	gI, gOk := query[guessKey]
	sI, sOk := query[scoreKey]
	if !gOk || !sOk {
		return nil, nil
	}
	guessSingle := gI[0]
	scoreSingle := sI[0]

	g := guess.New(guessSingle)
	s := score.New(scoreSingle)

	if len(guessSingle) == 0 && len(scoreSingle) == 0 {
		return nil, nil
	}
	var anyWord words.Words
	if err := g.Validate(anyWord); err != nil {
		return nil, fmt.Errorf("reading guess: %w", err)
	}

	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("reading score: %w", err)
	}

	r := result.Result{
		Guess: g,
		Score: s,
	}
	return &r, nil
}

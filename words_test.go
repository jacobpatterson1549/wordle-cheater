package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestNewWords(t *testing.T) {
	tests := []struct {
		input   string
		want    *words
		wantErr bool
	}{
		{
			want: &words{}, // no words
		},
		{
			input:   "tiny", // too short
			wantErr: true,
		},
		{
			input: "apple\nberry",
			want:  &words{"apple": {}, "berry": {}},
		},
		{
			input:   "APPLE", // uppercase
			wantErr: true,
		},
		{
			input: "extra\nbreak\nvalid\n\n",
			want:  &words{"extra": {}, "break": {}, "valid": {}},
		},
	}
	for i, test := range tests {
		got, gotErr := newWords(test.input)
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error creating words", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error creating words: %v", i, gotErr)
		case !reflect.DeepEqual(test.want, got):
			t.Errorf("test %v: words not equal:\nwanted: %v\ngot:    %v", i, test.want, got)
		}
	}
}

func TestWordsSorted(t *testing.T) {
	words := words{
		"abbey": {},
		"weary": {},
		"gravy": {},
	}
	if want, got := "abbey,gravy,weary", words.sorted(); want != got {
		t.Errorf("sorted words not equal:\nwanted: %q\ngot:    %q", want, got)
	}
}

func TestWordsCopy(t *testing.T) {
	w := "magic"
	a := words{
		w: {},
	}
	b := a.copy()
	if !reflect.DeepEqual(a, *b) {
		t.Errorf("copied values should be equal:\nwanted: %v\ngot:    %v", a, b)
	}
	delete(a, w)
	if reflect.DeepEqual(a, *b) {
		t.Errorf("copy should be to a different location in memory, both were %p", &a)
	}
}

func TestWordsScanShowPossible(t *testing.T) {
	tests := []struct {
		words
		in      string
		wantOut string
		wantErr bool
	}{
		{
			wantErr: true, // input EOF
		},
		{
			in:      "n",
			wantOut: "show possible words [Yn]: ",
		},
		{
			words:   words{"apple": {}, "berry": {}, "cakes": {}},
			in:      "NAH", // uppercase n is ok
			wantOut: "show possible words [Yn]: ",
		},
		{
			words:   words{"apple": {}, "berry": {}, "cakes": {}},
			in:      "yes\n",
			wantOut: "show possible words [Yn]: remaining valid words: apple,berry,cakes\n",
		},
		{
			words:   words{"apple": {}},
			in:      "\n", // user presses enter key (choosing default: Y)
			wantOut: "show possible words [Yn]: remaining valid words: apple\n",
		},
		{
			words:   words{"apple": {}},
			in:      "hmmm... no", // first word must be no
			wantOut: "show possible words [Yn]: remaining valid words: apple\n",
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.in)),
			Writer: bufio.NewWriter(&buf),
		}
		gotErr := test.words.scanShowPossible(rw)
		rw.Flush()
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error scanning to show possible words", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error scanning to show possible words: %v", i, gotErr)
		case test.wantOut != buf.String():
			t.Errorf("test %v: outputs not equal for possible words:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())
		}
	}
}

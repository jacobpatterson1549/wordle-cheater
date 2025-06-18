package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestNewScore(t *testing.T) {
	tests := []struct {
		in      string
		wantOut string
		want    score
		wantErr bool
	}{
		{
			wantErr: true, // input EOF
		},
		{
			in:      "ccccc",
			wantOut: "Enter score: ",
			want:    "ccccc",
		},
		{
			in:      "cCcCc",
			wantOut: "Enter score: ",
			want:    "ccccc",
		},
		{
			in:      "nac apple canac",
			wantOut: "Enter score: score must be 5 letters long\nEnter score: must be only the following letters: C, A, N\nEnter score: ",
			want:    "canac",
		},
	}
	for i, test := range tests {
		var buf strings.Builder
		rw := bufio.ReadWriter{
			Reader: bufio.NewReader(strings.NewReader(test.in)),
			Writer: bufio.NewWriter(&buf),
		}
		got, gotErr := newScore(rw)
		rw.Flush()
		switch {
		case test.wantErr:
			if gotErr == nil {
				t.Errorf("test %v: wanted error creating score", i)
			}
		case gotErr != nil:
			t.Errorf("test %v: unwanted error creating score: %v", i, gotErr)
		case test.want != *got:
			t.Errorf("test %v: scores not equal:\nwanted: %v\ngot:    %v", i, test.want, *got)
		case test.wantOut != buf.String():
			t.Errorf("test %v: new score outputs not equal:\nwanted: %q\ngot:    %q", i, test.wantOut, buf.String())
		}
	}
}

func TestScoreAllCorrect(t *testing.T) {
	tests := []struct {
		score
		want bool
	}{
		{
			score: "cccca",
		},
		{
			score: "c",
		},
		{
			score: "ccccc",
			want:  true,
		},
	}
	for i, test := range tests {
		if want, got := test.want, test.score == allCorrect; want != got {
			t.Errorf("test %v: allCorrect values not equal for %q: wanted %v, got %v", i, test.score, want, got)
		}
	}
}

package config

import (
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

var eh = flag.ContinueOnError

func TestNewConfigHelp(t *testing.T) {
	var sb strings.Builder
	_, err := newConfig(eh, &sb, "progrAm", "-h")
	if want := flag.ErrHelp; !errors.Is(err, want) {
		t.Fatalf("unwanted error: %v \n wanted: %v", err, want)
	}
	if want, got := "runs progrAm", sb.String(); !strings.Contains(got, want) {
		t.Errorf("wanted program help to contain %q \n got: %q", want, got)
	}
}

func TestNewConfigNoProgramName(t *testing.T) {
	if _, err := newConfig(eh, io.Discard); err == nil {
		t.Errorf("wanted error parsing args without program name")
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		env    [][]string
		wantOk bool
		want   Config
	}{
		{
			name:   "defaults",
			wantOk: true,
			want: Config{
				Port: "8000",
			},
		},
		{
			name: "all args",
			args: []string{
				"-port=1",
			},
			wantOk: true,
			want: Config{
				Port: "1",
			},
		},
		{
			name: "all env",
			args: []string{
				"-port=999", // environment wins
			},
			env: [][]string{
				{"PORT", "1"},
			},
			wantOk: true,
			want: Config{
				Port: "1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, e := range test.env {
				k, v := e[0], e[1]
				t.Setenv(k, v)
			}
			args := append([]string{"name"}, test.args...)
			got, err := newConfig(eh, io.Discard, args...)
			switch {
			case !test.wantOk:
				if err == nil {
					t.Errorf("wanted error")
				}
			case err != nil:
				t.Errorf("unwanted error: %v", err)
			default:
				got.fs = test.want.fs
				if test.want != *got {
					t.Errorf("not equal: \n wanted: %#v \n got:    %#v", test.want, *got)
				}
			}
		})
	}
}

package server

import (
	"fmt"
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantCode int
	}{
		{
			name:     "unknown",
			target:   "/unknown",
			wantCode: 404,
		},
		{
			name:     "wordle-empty",
			target:   wordlePath,
			wantCode: 200,
		},
		{
			name:     "wordle-ok",
			target:   wordlePath + "?g0=words&s0=ccccc",
			wantCode: 200,
		},
		{
			name:     "wordle-bad",
			target:   wordlePath + "?g0=word&s0=c",
			wantCode: 400,
		},
		{
			name:     "spelling-bee-empty",
			target:   spellingBeePath,
			wantCode: 200,
		},
		{
			name:     "spelling-bee-ok",
			target:   spellingBeePath + "?" + centralLetterParam + "=a&" + otherLettersParam + "=bcdefg",
			wantCode: 200,
		},
		{
			name:     "spelling-bee-bad",
			target:   spellingBeePath + "?" + centralLetterParam + "=az&" + otherLettersParam + "=bcdefgh",
			wantCode: 400,
		},
		{
			name:     "spelling-bee-404",
			target:   spellingBeePath + "-404",
			wantCode: 404,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var wordsText string
			h := NewHandler(wordsText)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", test.target, nil)
			h.ServeHTTP(w, r)
			if want, got := test.wantCode, w.Result().StatusCode; want != got {
				t.Errorf("wanted %v, got %v (body: %q)", want, got, w.Body.String())
			}
		})
	}
}

func TestHandleBadRequest(t *testing.T) {
	message := "my-message"
	err := fmt.Errorf("my-error")
	w := httptest.NewRecorder()
	handleBadRequest(w, message, err)
	if want, got := 400, w.Result().StatusCode; want != got {
		t.Errorf("status codes: wanted %v, got %v", want, got)
	}
	body := strings.TrimSpace(w.Body.String())
	if want, got := "my-message: my-error", body; want != got {
		t.Errorf("response bodies: wanted: %q got: %q", want, got)
	}
}

func TestHandleTemplateError(t *testing.T) {
	tests := []struct {
		name          string
		data          any
		wantBodyStart string
		wantCode      int
	}{
		{"ok", struct{ A int }{A: 9}, "9", 200},
		{"bad", nil, "rendering", 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			h := Handler{
				tmpl: template.Must(newTemplate().Parse("{{.A}}")),
			}
			w := httptest.NewRecorder()
			h.handleTemplate(w, test.data)
			if want, got := test.wantCode, w.Result().StatusCode; want != got {
				t.Errorf("status codes: wanted %v, got %v", want, got)
			}
			got := w.Body.String()
			if !strings.HasPrefix(got, test.wantBodyStart) {
				t.Errorf("response body: wanted prefix: %q got: %q", test.wantBodyStart, got)
			}
		})
	}
}

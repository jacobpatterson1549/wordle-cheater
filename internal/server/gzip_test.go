package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithContentEncoding(t *testing.T) {
	msg := "OK_gzip"
	tests := []struct {
		name    string
		ae      string
		wantCE  string
		getBody func(t *testing.T, r io.Reader) io.Reader
	}{
		{
			name:   "gzip",
			ae:     "gzip, deflate, br",
			wantCE: "gzip",
			getBody: func(t *testing.T, r io.Reader) io.Reader {
				t.Helper()
				gr, err := gzip.NewReader(r)
				if err != nil {
					t.Fatalf("creating gzip reader: %v", err)
				}
				return gr
			},
		},
		{
			name:   "UNKNOWN",
			ae:     "UNKNOWN",
			wantCE: "",
			getBody: func(t *testing.T, r io.Reader) io.Reader {
				return r
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h1 := func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(msg))
			}
			h2 := withContentEncoding(http.HandlerFunc(h1))
			w := httptest.NewRecorder()
			r := httptest.NewRequest("", "/", nil)
			r.Header.Add("Accept-Encoding", test.ae)
			h2.ServeHTTP(w, r)
			gotHeader := w.Header()
			gotCE := gotHeader.Get("Content-Encoding")
			if test.wantCE != gotCE {
				t.Fatalf("wanted %q Content-Encoding, got: %q",
					test.wantCE, gotCE)
			}
			gr := test.getBody(t, w.Body)
			b, err := io.ReadAll(gr)
			if err != nil {
				t.Fatalf("reading gzip encoded message: %v", err)
			}
			if want, got := msg, string(b); want != got {
				t.Errorf("body not encoded as desired: wanted %q, got %q",
					want, got)
			}
		})
	}
}

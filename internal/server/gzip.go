package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func withContentEncoding(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enc := r.Header.Get("Accept-Encoding")
		if strings.Contains(enc, "gzip") {
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			wrw := wrappedResponseWriter{
				Writer:         gzw,
				ResponseWriter: w,
			}
			wrw.Header().Set("Content-Encoding", "gzip")
			w = wrw
		}
		h.ServeHTTP(w, r)
	}
}

type wrappedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (wrw wrappedResponseWriter) Write(p []byte) (n int, err error) {
	return wrw.Writer.Write(p)
}

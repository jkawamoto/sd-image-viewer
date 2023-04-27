package server

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	Code int
}

func (w *responseWriter) WriteHeader(code int) {
	w.Code = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogger(h http.Handler, logger *log.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		w := &responseWriter{
			ResponseWriter: res,
		}
		h.ServeHTTP(w, req)
		logger.Printf("%v %v %v", w.Code, req.Method, req.URL)
	}
}

package middleware

import (
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriter{ResponseWriter: w, StatusCode: 200}
		next.ServeHTTP(rw, r)
		fmt.Printf("[%s] %s - %d\n", r.Method, r.URL.Path, rw.StatusCode)
	})
}
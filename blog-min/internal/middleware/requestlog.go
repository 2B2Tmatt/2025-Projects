package middleware

import (
	"log"
	"net/http"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s Host=%s Cookie=%t", r.Method, r.URL.Path, r.Host, r.Header.Get("Cookie") != "")
		next.ServeHTTP(w, r)
	})
}

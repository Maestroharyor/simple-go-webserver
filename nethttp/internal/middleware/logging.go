package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusOrder struct {
	http.ResponseWriter
	status int
}

func (s *statusOrder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusOrder{
			ResponseWriter: w, status: http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.status, duration)
	})
}

package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ycatan-w/himawari-api/internal/output"
	"github.com/ycatan-w/himawari-api/internal/output/colors"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	return lrw.ResponseWriter.Write(b)
}

func formatMethod(method string) string {
	return fmt.Sprintf("%-20s", fmt.Sprintf("[%s]", colors.Yellow(method)))
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s][Request ]%s → %s", output.AppNameGreen(), formatMethod(r.Method), r.URL.Path)
		start := time.Now()

		// override the response writer to intercept the response code
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		status := lrw.statusCode
		if status == 0 {
			status = http.StatusOK
		}
		log.Printf("[%s][Response]%s ← %s ∙ %d (%v)", output.AppNameGreen(), formatMethod(r.Method), r.URL.Path, status, duration)
	})
}

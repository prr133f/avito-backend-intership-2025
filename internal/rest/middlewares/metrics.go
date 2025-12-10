package middlewares

import (
	"net/http"

	"github.com/prr133f/avito-backend-intership-2025/internal/metrics"
)

func CountRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestsTotal.WithLabelValues(r.URL.Path).Inc()
		next.ServeHTTP(w, r)
	})
}

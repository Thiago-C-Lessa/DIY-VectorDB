package custom_middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(requests *prometheus.CounterVec, duration *prometheus.HistogramVec) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rec := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rec, r)

			elapsed := time.Since(start).Seconds()

			path := chi.RouteContext(r.Context()).RoutePattern()
			method := r.Method
			status := strconv.Itoa(rec.status)

			requests.WithLabelValues(path, method, status).Inc()
			duration.WithLabelValues(path, method, status).Observe(elapsed)
		})
	}
}

package main

import (
	"DIY-VectorDB/internal/db"
	"DIY-VectorDB/internal/http/custom_middleware"
	"DIY-VectorDB/internal/http/server"

	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of car requests.",
		}, []string{"path", "method", "status"})

	resquestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:                       "",
			Subsystem:                       "",
			Name:                            "http_request_duration",
			Help:                            "duration of http requests",
			ConstLabels:                     nil,
			Buckets:                         nil,
			NativeHistogramBucketFactor:     0,
			NativeHistogramZeroThreshold:    0,
			NativeHistogramMaxBucketNumber:  0,
			NativeHistogramMinResetDuration: 0,
			NativeHistogramMaxZeroThreshold: 0,
			NativeHistogramMaxExemplars:     0,
			NativeHistogramExemplarTTL:      0,
		}, []string{"path", "method", "status"})
)

func initPrometheus() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(resquestDuration)
}

func main() {
	initPrometheus()

	dbvm := db.NewVecMemDB()

	r := server.NewRouter()

	r.Use(custom_middleware.MetricsMiddleware(httpRequestsTotal, resquestDuration))

	r.Handle("/metrics", promhttp.Handler())

	r.Mount("/store", server.StoreRoutes(dbvm))
	r.Mount("/fetch", server.FetchRoutes(dbvm))
	r.Mount("/update", server.UpdateRoutes(dbvm))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

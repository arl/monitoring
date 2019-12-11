package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	mux   *http.ServeMux
	cache Cache
}

func newServer(csize int) *server {
	return &server{
		mux:   http.NewServeMux(),
		cache: NewLRUCache(csize),
	}
}

func (s *server) handleAdd(w http.ResponseWriter, r *http.Request) {
	// Extract the (key, value) pair to add to the cache
	query := r.URL.Query()
	k, v := query.Get("k"), query.Get("v")
	if k == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	s.cache.Add(k, v)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	// Extract the key to lookup in cache
	query := r.URL.Query()
	k := query.Get("k")
	if k == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Cache lookup
	v, ok := s.cache.Get(k)
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		cacheMisses.Inc()
	} else {
		cacheHits.Inc()
	}

	fmt.Fprint(w, v)
}

func (s *server) serve(addr string) error {
	log.Println("server starting:", addr)
	return http.ListenAndServe(addr, s.mux)
}

func recordMetrics(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		totalRequests.Inc()

		// Measure taken by the handler h
		t0 := time.Now()
		h(w, r)
		duration := time.Since(t0) / time.Microsecond
		requestDuration.Observe(float64(duration))
	}
}

// IMPORT
// 	"github.com/prometheus/client_golang/prometheus/promhttp"

func (s *server) setupRoutes() {
	s.mux.HandleFunc("/add", recordMetrics(s.handleAdd))
	s.mux.HandleFunc("/get", recordMetrics(s.handleGet))
	s.mux.Handle("/metrics", promhttp.Handler())
}

var (
	totalRequests = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_requests_total",
			Help: "The total number of requests",
		})

	cacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "The total number of cache hits",
		})

	cacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "The total number of cache misses",
		})

	requestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_duration_microseconds",
			Help:    "The duration of requests",
			Buckets: prometheus.LinearBuckets(0, 5, 20),
		})
)

func main() {
	addr := flag.String("addr", ":8080", "server listen address")
	csize := flag.Int("size", 256, "LRU cache size")

	flag.Parse()

	s := newServer(*csize)
	s.setupRoutes()

	log.Fatal(s.serve(*addr))
}

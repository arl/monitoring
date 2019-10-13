package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	}

	fmt.Fprint(w, v)
}

func (s *server) serve(addr string) error {
	log.Println("server starting:", addr)
	return http.ListenAndServe(":8080", s.mux)
}

func (s *server) setupRoutes() {
	s.mux.HandleFunc("/add", s.handleAdd)
	s.mux.HandleFunc("/get", s.handleGet)
}

func main() {
	addr := flag.String("addr", ":8080", "server listen address")
	csize := flag.Int("size", 256, "LRU cache size")

	flag.Parse()

	s := newServer(*csize)
	s.setupRoutes()

	log.Fatal(s.serve(*addr))
}

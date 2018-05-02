package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/zgiber/cache"
)

const maxBytesLimit = 1024 * 1024 * 128 // 128 MB
const defaultTTL = 168 * time.Hour      // One week

type service struct {
	c cache.Cache
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "":
		s.handleGet(w, r)
	case "PUT", "POST":
		s.handleSet(w, r)
	case "DELETE":
		s.handleDelete(w, r)
	}
}

func newService() *service {
	return &service{
		c: cache.NewMemCache(cache.MaxBytesLimit(maxBytesLimit)),
	}
}

func (s *service) handleSet(w http.ResponseWriter, r *http.Request) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = s.c.Set(r.URL.Path, value, defaultTTL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *service) handleGet(w http.ResponseWriter, r *http.Request) {
	value, err := s.c.Fetch(r.URL.Path)
	if err != nil {
		if err == cache.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *service) handleDelete(w http.ResponseWriter, r *http.Request) {
	s.c.Delete(r.URL.Path)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	if err := http.ListenAndServe(":8080", newService()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

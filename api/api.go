package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gavincabbage/kurz/store"
)

type Server struct {
	store LinkStore
}

type LinkStore interface {
	Put(context.Context, string, string) error
	Get(context.Context, string) (string, error)
}

func New(logger *log.Logger, address string, store LinkStore) *http.Server {
	api := Server{store}

	router := http.NewServeMux()
	router.Handle("/", &api)
	router.HandleFunc("/health", health)

	return &http.Server{
		Addr:         address,
		Handler:      logging(logger)(router), // TODO(gavincabbage): Make logger optional.
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.postLink(w, r)
	case http.MethodGet:
		s.getLink(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *Server) postLink(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("reading request body: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := url.Parse(string(link)); err != nil {
		http.Error(w, fmt.Sprintf("invalid link %q: %v", link, err), http.StatusBadRequest)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if len(key) == 0 {
		http.Error(w, "missing link key", http.StatusBadRequest)
		return
	}

	if err := s.store.Put(r.Context(), key, string(link)); err != nil {
		http.Error(w, fmt.Sprintf("storing link: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getLink(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")
	if len(key) == 0 {
		http.Error(w, "missing link key", http.StatusBadRequest)
		return
	}

	v, err := s.store.Get(r.Context(), key)
	if err != nil {
		notFound := store.NotFound(key)
		if errors.As(err, &notFound) {
			http.Error(w, fmt.Sprintf("not found: %s", notFound), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("getting link: %v", err), http.StatusInternalServerError)
		}
		return

	}

	http.Redirect(w, r, v, http.StatusPermanentRedirect)
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(),
					"Content-Type:", r.Header.Get("content-type"),
					"Accept:", r.Header.Get("accept"))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

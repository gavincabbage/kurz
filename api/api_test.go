package api_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gavincabbage/kurz/api"
	"github.com/gavincabbage/kurz/store/inmem"
)

func TestServer(t *testing.T) {
	subject := api.New(log.Default(), ":8080", &inmem.Store{})
	url := "http://test.test"
	{
		r := httptest.NewRequest(http.MethodPost, "/foo", strings.NewReader(url))
		w := httptest.NewRecorder()

		subject.Handler.ServeHTTP(w, r)

		response := w.Result()
		defer response.Body.Close()

		if want, got := http.StatusCreated, response.StatusCode; want != got {
			t.Errorf("wrong status: expected %v got %v", want, got)
		}
	}
	{
		r := httptest.NewRequest(http.MethodGet, "/foo", nil)
		w := httptest.NewRecorder()
		subject.Handler.ServeHTTP(w, r)

		response := w.Result()
		defer response.Body.Close()

		if want, got := http.StatusPermanentRedirect, response.StatusCode; want != got {
			t.Errorf("wrong status: expected %v got %v", want, got)
		}

		got, _ := io.ReadAll(response.Body)
		if !strings.Contains(string(got), url) {
			t.Errorf("expected response to redirect to %v but got body %q", url, got)
		}
	}
	{
		r := httptest.NewRequest(http.MethodGet, "/dne", nil)
		w := httptest.NewRecorder()
		subject.Handler.ServeHTTP(w, r)

		response := w.Result()
		defer response.Body.Close()

		if want, got := http.StatusNotFound, response.StatusCode; want != got {
			t.Errorf("wrong status: expected %v got %v", want, got)
		}
	}
}

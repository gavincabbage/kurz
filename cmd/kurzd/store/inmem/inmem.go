package inmem

import (
	"sync"

	"github.com/gavincabbage/kurz/cmd/kurzd/store"
)

type Store struct {
	kv sync.Map
}

func (s *Store) Put(k, v string) error {
	s.kv.Store(k, v)
	return nil
}

func (s *Store) Get(k string) (string, error) {
	v, ok := s.kv.Load(k)
	if !ok {
		return "", store.NotFound(k)
	}
	return v.(string), nil
}

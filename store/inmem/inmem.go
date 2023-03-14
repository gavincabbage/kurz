package inmem

import (
	"context"
	"sync"

	"github.com/gavincabbage/kurz/store"
)

type Store struct {
	kv sync.Map
}

func (s *Store) Put(_ context.Context, k, v string) error {
	s.kv.Store(k, v)
	return nil
}

func (s *Store) Get(_ context.Context, k string) (string, error) {
	v, ok := s.kv.Load(k)
	if !ok {
		return "", store.NotFound(k)
	}
	return v.(string), nil
}

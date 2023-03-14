package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/gavincabbage/kurz/store"
)

type Store struct {
	conn *redis.Client
}

func New(conn *redis.Client) *Store {
	return &Store{conn}
}

func (s *Store) Put(ctx context.Context, k, v string) error {
	if err := s.conn.Set(ctx, k, v, 0); err != nil {
		return fmt.Errorf("setting key %s to value %s: %v", k, v, err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, k string) (string, error) {
	v, err := s.conn.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", store.NotFound(k)
		}
		return "", fmt.Errorf("getting key %s: %v", k, err)
	}
	return v, nil
}

package store

import "fmt"

type ErrNotFound string

func NotFound(key string) error {
	return ErrNotFound(key)
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("key not found: %s", e)
}

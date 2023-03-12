package store

import "fmt"

type ErrNotFound struct {
	key string
}

func NotFound(key string) *ErrNotFound {
	return &ErrNotFound{key}
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("key not found: %s", e.key)
}

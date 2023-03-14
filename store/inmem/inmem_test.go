package inmem_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gavincabbage/kurz/store"
	"github.com/gavincabbage/kurz/store/inmem"
)

func TestStore(t *testing.T) {
	subject := &inmem.Store{}
	{
		if err := subject.Put(context.Background(), "foo", "bar"); err != nil {
			t.Errorf("failed to put key initially")
		}
		if got, err := subject.Get(context.Background(), "foo"); err != nil {
			t.Errorf("failed to get key initially")
		} else if got != "bar" {
			t.Errorf("wrong value; expected \"bar\" but got %v", got)
		}
	}
	{
		if err := subject.Put(context.Background(), "foo", "baz"); err != nil {
			t.Errorf("failed to put key initially")
		}
		if got, err := subject.Get(context.Background(), "foo"); err != nil {
			t.Errorf("failed to get key initially")
		} else if got != "baz" {
			t.Errorf("wrong value; expected \"baz\" but got %v", got)
		}
	}
	{
		want := store.NotFound("dne")
		if _, err := subject.Get(context.Background(), "dne"); !errors.As(err, &want) {
			t.Errorf("wrong error; expected ErrNotFound but got %v", err)
		}
	}
}

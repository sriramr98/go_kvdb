package store

import (
	"fmt"
	"testing"
)

func TestInMemStore(t *testing.T) {
	store := NewInMemoryStore()

	t.Run("get key after put returns value", func(t *testing.T) {
		store.Set("foo", NewDataObj([]byte("bar"), 0))
		value, err := store.Get("foo")
		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}

		if string(value.value) != "bar" {
			t.Errorf("Expected value to be 'bar', got %v", value)
		}
	})

	t.Run("get key that doesn't exist returns error", func(t *testing.T) {
		_, err := store.Get("test1")

		if err != ErrKeyNotFound {
			t.Errorf("Expected err to be ErrKeyNotFound, got %v", err)
		}
	})

	t.Run("delete key removes key", func(t *testing.T) {
		store.Set("test2", NewDataObj([]byte("value2"), 0))
		store.Delete("test2")
		_, err := store.Get("test2")

		if err != ErrKeyNotFound {
			t.Errorf("Expected err to be ErrKeyNotFound, got %v", err)
		}
	})
}

func BenchmarkInMemStore(b *testing.B) {
	b.Run("set", func(b *testing.B) {
		store := NewInMemoryStore()

		for i := 0; i < b.N; i++ {
			store.Set(fmt.Sprintf("foo%q", i), NewDataObj([]byte("bar"), 0))
		}
	})

	b.Run("get", func(b *testing.B) {
		store := NewInMemoryStore()

		for i := 0; i < b.N; i++ {
			store.Get(fmt.Sprintf("foo%q", i))
		}
	})

	b.Run("delete", func(b *testing.B) {
		store := NewInMemoryStore()

		for i := 0; i < b.N; i++ {
			store.Delete(fmt.Sprintf("foo%q", i))
		}
	})
}

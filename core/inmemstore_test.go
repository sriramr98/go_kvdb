package core

import (
	"fmt"
	"testing"
)

func TestInMemStore(t *testing.T) {

	store := NewInMemoryStore()

	t.Run("get key after put returns value", func(t *testing.T) {
		store.Set("foo", []byte("bar"))
		value, err := store.Get("foo")

		if err != nil {
			t.Errorf("Expected err to be nil, got %v", err)
		}

		if string(value) != "bar" {
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
		store.Set("test2", []byte("foo"))
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
			store.Set(fmt.Sprintf("foo%q", i), []byte("bar"))
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

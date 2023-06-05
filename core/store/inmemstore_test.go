package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemStore(t *testing.T) {

	store := NewInMemoryStore()

	t.Run("get key after put returns value", func(t *testing.T) {
		store.Set("foo", []byte("bar"))
		value, err := store.Get("foo")

		assert.NoError(t, err)
		assert.Equal(t, []byte("bar"), value)
	})

	t.Run("get key that doesn't exist returns error", func(t *testing.T) {
		_, err := store.Get("test1")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})

	t.Run("delete key removes key", func(t *testing.T) {
		store.Set("test2", []byte("foo"))
		store.Delete("test2")
		_, err := store.Get("test2")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyNotFound)
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

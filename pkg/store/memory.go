package store

import (
	"context"
	"sync"
)

type Memory[T any] struct {
	mu    sync.RWMutex
	store map[string]T
}

func NewMemoryStore[T any]() *Memory[T] {
	return &Memory[T]{
		store: make(map[string]T),
	}
}

func (s *Memory[T]) Get(_ context.Context, key string) (T, error) {
	return s.store[key], nil
}

func (s *Memory[T]) Delete(_ context.Context, key string) error {
	delete(s.store, key)
	return nil
}

func (s *Memory[T]) Set(_ context.Context, key string, value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = value
	return nil
}

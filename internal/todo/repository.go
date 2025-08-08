package todo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("todo not found")
)

type Repository interface {
	Create(ctx context.Context, title string) (Todo, error)
	Get(ctx context.Context, id string) (Todo, error)
	List(ctx context.Context) ([]Todo, error)
	Update(ctx context.Context, id string, update UpdateTodoRequest) (Todo, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]Todo
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{store: make(map[string]Todo)}
}

func (r *InMemoryRepository) Create(ctx context.Context, title string) (Todo, error) {
	_ = ctx
	now := time.Now().UTC()
	t := Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.mu.Lock()
	r.store[t.ID] = t
	r.mu.Unlock()
	return t, nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (Todo, error) {
	_ = ctx
	r.mu.RLock()
	t, ok := r.store[id]
	r.mu.RUnlock()
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

func (r *InMemoryRepository) List(ctx context.Context) ([]Todo, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Todo, 0, len(r.store))
	for _, t := range r.store {
		result = append(result, t)
	}
	return result, nil
}

func (r *InMemoryRepository) Update(ctx context.Context, id string, update UpdateTodoRequest) (Todo, error) {
	_ = ctx
	r.mu.Lock()
	t, ok := r.store[id]
	if !ok {
		r.mu.Unlock()
		return Todo{}, ErrNotFound
	}
	if update.Title != nil {
		t.Title = *update.Title
	}
	if update.Completed != nil {
		t.Completed = *update.Completed
	}
	t.UpdatedAt = time.Now().UTC()
	r.store[id] = t
	r.mu.Unlock()
	return t, nil
}

func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return ErrNotFound
	}
	delete(r.store, id)
	return nil
}

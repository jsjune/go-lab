package store

import (
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("item not found")

type Item struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// 동시 요청을 위해 RWMutex로 보호 (읽기는 동시에, 쓰기는 단독으로)
type InMemoryStore struct {
	mu    sync.RWMutex
	items map[int64]Item
	next  int64
}

func New() *InMemoryStore {
	return &InMemoryStore{items: make(map[int64]Item), next: 1}
}

func (s *InMemoryStore) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

func (s *InMemoryStore) Get(id int64) (*Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &item, nil
}

func (s *InMemoryStore) Create(name, description string) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	item := Item{ID: s.next, Name: name, Description: description, CreatedAt: time.Now()}
	s.items[s.next] = item
	s.next++
	return item
}

func (s *InMemoryStore) Update(id int64, name, description string) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	item.Name = name
	item.Description = description
	s.items[id] = item
	return &item, nil
}

func (s *InMemoryStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

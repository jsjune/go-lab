package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var ErrNotFound = errors.New("item not found")

type Item struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	path  string
	Items []Item `json:"items"`
	Next  int64  `json:"next_id"`
}

func New() (*Store, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ".item-cli", "data.json")

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	s := &Store{path: path, Next: 1}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

func (s *Store) List() []Item {
	return s.Items
}

func (s *Store) Get(id int64) (*Item, error) {
	for _, item := range s.Items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, ErrNotFound
}

func (s *Store) Create(name, description string) (*Item, error) {
	item := Item{
		ID:          s.Next,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	s.Items = append(s.Items, item)
	s.Next++
	return &item, s.save()
}

func (s *Store) Delete(id int64) error {
	for i, item := range s.Items {
		if item.ID == id {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			return s.save()
		}
	}
	return ErrNotFound
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

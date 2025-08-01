// File: internal/task/store.go
package task

import (
	"fmt"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
	Clock int    `json:"clock"` // Add Clock field for Lamport timestamp
}

// Store interface is correct and already includes Get
type Store interface {
	List() ([]Task, error)
	Create(t Task) (Task, error)
	Update(id int, t Task) (Task, error)
	Delete(id int) error
	Get(id int) (Task, error)
}

// InMemoryStore is updated to use int IDs
type InMemoryStore struct {
	mu    sync.Mutex
	tasks map[int]Task
	next  int
}

func NewInMemoryStore() *InMemoryStore {
	// Initialize with next=0 so the first ID is 1
	return &InMemoryStore{tasks: make(map[int]Task), next: 0}
}

// Implement the missing Get method for InMemoryStore
func (s *InMemoryStore) Get(id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task with id %d not found", id)
	}
	return t, nil
}

func (s *InMemoryStore) List() ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out, nil
}

func (s *InMemoryStore) Create(t Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	t.ID = s.next
	s.tasks[t.ID] = t
	return t, nil
}

func (s *InMemoryStore) Update(id int, t Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task with id %d not found", id)
	}
	t.ID = id
	s.tasks[id] = t
	return t, nil
}

func (s *InMemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task with id %d not found", id)
	}
	delete(s.tasks, id)
	return nil
}
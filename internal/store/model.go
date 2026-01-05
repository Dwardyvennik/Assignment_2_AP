package store

import (
	"sync"
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID        string     `json:"id"`
	Payload   string     `json:"payload"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type Stats struct {
	mu         sync.RWMutex
	Submitted  int `json:"submitted"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) IncrementSubmitted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Submitted++
}

func (s *Stats) IncrementCompleted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Completed++
}

func (s *Stats) SetInProgress(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InProgress = count
}

// GetStats copy  current statistics
func (s *Stats) GetStats() (int, int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Submitted, s.Completed, s.InProgress
}

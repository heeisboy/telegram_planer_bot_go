package storage

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// MemoryStore keeps tasks in memory with a per-user map.
type MemoryStore struct {
	mu    sync.Mutex
	users map[int64]*userBucket
}

type userBucket struct {
	nextID int
	tasks  map[int]Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[int64]*userBucket),
	}
}

func (s *MemoryStore) CreateTask(userID int64, text string) (Task, error) {
	if text == "" {
		return Task{}, errors.New("text is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	bucket := s.ensureUser(userID)
	bucket.nextID++
	now := time.Now().UTC()
	task := Task{
		ID:        bucket.nextID,
		UserID:    userID,
		Text:      text,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	bucket.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryStore) ListTasks(userID int64) ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, ok := s.users[userID]
	if !ok {
		return []Task{}, nil
	}

	out := make([]Task, 0, len(bucket.tasks))
	for _, task := range bucket.tasks {
		out = append(out, task)
	}
	// Stable output for consistent UX/tests.
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *MemoryStore) GetTask(userID int64, id int) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, ok := s.users[userID]
	if !ok {
		return Task{}, false, nil
	}
	task, ok := bucket.tasks[id]
	return task, ok, nil
}

func (s *MemoryStore) UpdateTask(userID int64, task Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, ok := s.users[userID]
	if !ok {
		return fmt.Errorf("user %d: %w", userID, ErrUserNotFound)
	}
	if _, ok := bucket.tasks[task.ID]; !ok {
		return fmt.Errorf("task %d: %w", task.ID, ErrTaskNotFound)
	}
	task.UpdatedAt = time.Now().UTC()
	bucket.tasks[task.ID] = task
	return nil
}

func (s *MemoryStore) DeleteTask(userID int64, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket, ok := s.users[userID]
	if !ok {
		return fmt.Errorf("user %d: %w", userID, ErrUserNotFound)
	}
	if _, ok := bucket.tasks[id]; !ok {
		return fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
	}
	delete(bucket.tasks, id)
	return nil
}

// ListUserIDs returns known user IDs for scheduling.
func (s *MemoryStore) ListUserIDs() []int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]int64, 0, len(s.users))
	for userID := range s.users {
		out = append(out, userID)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (s *MemoryStore) ensureUser(userID int64) *userBucket {
	bucket, ok := s.users[userID]
	if !ok {
		bucket = &userBucket{nextID: 0, tasks: make(map[int]Task)}
		s.users[userID] = bucket
	}
	return bucket
}

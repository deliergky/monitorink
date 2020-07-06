package store

import (
	"context"
	"sync"

	"github.com/deliergky/monitorink/data"
)

// InMemoryStore implementation of the store interface
type InMemoryStore struct {
	sync.Mutex
	messages []data.ResponseData
}

func (s *InMemoryStore) Persist(_ context.Context, d data.ResponseData) error {
	s.Lock()
	defer s.Unlock()
	if len(s.messages) == 0 {
		s.messages = make([]data.ResponseData, 0)
	}
	s.messages = append(s.messages, d)
	return nil
}

func (s *InMemoryStore) Retrieve(_ context.Context) ([]data.ResponseData, error) {
	return s.messages, nil
}

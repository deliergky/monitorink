package queue

import (
	"context"
	"errors"
	"sync"

	"github.com/deliergky/monitorink/data"
)

var ErrorEmpty = errors.New("empty_queue")

// InMemoryQueue implementation of the Queue interface
type InMemoryQueue struct {
	sync.Mutex
	messages []data.ResponseData
}

func (q *InMemoryQueue) Push(d data.ResponseData) error {
	q.Lock()
	defer q.Unlock()
	if len(q.messages) == 0 {
		q.messages = make([]data.ResponseData, 0)
	}
	q.messages = append(q.messages, d)
	return nil
}

func (q *InMemoryQueue) Pop() (data.ResponseData, error) {
	q.Lock()
	defer q.Unlock()
	if len(q.messages) == 0 {
		return data.ResponseData{}, ErrorEmpty
	}
	d := q.messages[len(q.messages)-1]
	q.messages = q.messages[:len(q.messages)-1]
	return d, nil
}

func (q *InMemoryQueue) Next(_ context.Context) bool {
	q.Lock()
	defer q.Unlock()
	return len(q.messages) > 0
}

func (q *InMemoryQueue) Current() data.ResponseData {
	d, _ := q.Pop()
	return d
}

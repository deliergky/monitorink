package monitor

import (
	"context"
	"log"

	"github.com/deliergky/monitorink/data"
	"github.com/deliergky/monitorink/pipeline"
	"github.com/deliergky/monitorink/queue"
	"github.com/deliergky/monitorink/store"
)

type HeartbeatConsumer struct {
	q queue.Queue
}

func NewHeartbeatConsumer(q queue.Queue) *HeartbeatConsumer {
	return &HeartbeatConsumer{
		q: q,
	}
}

func (c HeartbeatConsumer) Consume(ctx context.Context, d pipeline.Data) error {
	rd := d.(data.ResponseData)
	return c.q.Push(rd)
}

type ResponseConsumer struct {
	s store.Store
}

func (c ResponseConsumer) Consume(ctx context.Context, d pipeline.Data) error {
	rd := d.(data.ResponseData)
	log.Printf("Persisting %v", rd)
	return c.s.Persist(ctx, rd)
}

func NewResponseConsumer(s store.Store) *ResponseConsumer {
	return &ResponseConsumer{
		s: s,
	}
}

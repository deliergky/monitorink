package monitor

import (
	"context"
	"time"

	"github.com/deliergky/monitorink/pipeline"
	"github.com/deliergky/monitorink/queue"
)

type ResponseProducer struct {
	q queue.Queue
}

func (p ResponseProducer) Next(ctx context.Context) bool {
	return p.q.Next(ctx)
}

func (p ResponseProducer) Current() pipeline.Data {
	return p.q.Current()
}

func NewResponseProducer(q queue.Queue) *ResponseProducer {
	return &ResponseProducer{
		q: q,
	}
}

type HeartbeatProducer struct {
	ticker *time.Ticker
}

func (p HeartbeatProducer) Next(ctx context.Context) bool {
	select {
	case <-p.ticker.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (p HeartbeatProducer) Current() pipeline.Data {
	return nil
}

func NewHeartbeatProducer(interval time.Duration) *HeartbeatProducer {
	return &HeartbeatProducer{
		ticker: time.NewTicker(interval),
	}
}

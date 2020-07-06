package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/deliergky/monitorink/data"
	"github.com/deliergky/monitorink/pipeline"
)

type Heartbeat struct {
	config   HeartbeatConfig
	pipeline *pipeline.Pipeline
	producer pipeline.Producer
	consumer pipeline.Consumer
}

type HeartbeatConfig struct {
	URL       string
	SearchFor string
	Interval  time.Duration
}

func check(config HeartbeatConfig) pipeline.ProcessorFunc {
	return func(ctx context.Context, _ pipeline.Data) (pipeline.Data, error) {
		rd := data.ResponseData{
			URL: config.URL,
		}
		start := time.Now()
		r, err := http.Get(config.URL)
		if err != nil {
			return nil, err
		}
		rd.ResponseTime = time.Now().Sub(start).Milliseconds()
		defer r.Body.Close()
		rd.StatusCode = r.StatusCode
		return rd, nil
	}
}

func NewHeartbeat(config HeartbeatConfig, producer pipeline.Producer, consumer pipeline.Consumer) *Heartbeat {

	return &Heartbeat{
		config:   config,
		consumer: consumer,
		producer: producer,
		pipeline: pipeline.New(pipeline.FIFO(check(config))),
	}
}

func (h *Heartbeat) Beat(ctx context.Context) error {
	return h.pipeline.Process(ctx, h.producer, h.consumer)
}

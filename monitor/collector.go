package monitor

import (
	"context"

	"github.com/deliergky/monitorink/data"
	"github.com/deliergky/monitorink/pipeline"
)

type Collector struct {
	pipeline *pipeline.Pipeline
	producer pipeline.Producer
	consumer pipeline.Consumer
}

func collect() pipeline.ProcessorFunc {
	return func(ctx context.Context, in pipeline.Data) (pipeline.Data, error) {
		out := in.(data.ResponseData)
		return out, nil
	}
}

func NewCollector(producer pipeline.Producer, consumer pipeline.Consumer) *Collector {
	return &Collector{
		pipeline: pipeline.New(pipeline.FIFO(collect())),
		producer: producer,
		consumer: consumer,
	}
}

func (h *Collector) Collect(ctx context.Context) error {
	return h.pipeline.Process(ctx, h.producer, h.consumer)

}

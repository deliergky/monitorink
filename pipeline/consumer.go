package pipeline

import (
	"context"
)

// Consumer describes the behaviour of the final step/job of a pipeline
// Ex: Producing data.ResponseData messages on a Kafka Topic
type Consumer interface {
	Consume(context.Context, Data) error
}

package pipeline

import (
	"context"
)

// Producer is implemented by types that can produce data to feed the pipeline execution
type Producer interface {
	Next(context.Context) bool
	Current() Data
}

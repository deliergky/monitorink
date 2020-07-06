package pipeline

import "context"

// ProcessorFunc is a helper type that implements the Processor interface
type ProcessorFunc func(context.Context, Data) (Data, error)

// Processor is implemented by types that can process the data during pipeline execution
type Processor interface {
	Process(context.Context, Data) (Data, error)
}

// Process applies ProcessorFunc on the Data provided
func (f ProcessorFunc) Process(ctx context.Context, m Data) (Data, error) {
	return f(ctx, m)
}

package pipeline

import (
	"context"
	"log"
)

// StageResource represents the input and output data for a Stage
type StageResource interface {
	Input() <-chan Data
	Output() chan<- Data
	Error() chan<- error
}

// Stage represents one step in the pipeline
type Stage interface {
	Run(context.Context, StageResource)
}

// stageResource implement the StageResource interface
type stageResource struct {
	input  <-chan Data
	output chan<- Data
	err    chan<- error
}

func (rc *stageResource) Input() <-chan Data {
	return rc.input
}

func (rc *stageResource) Output() chan<- Data {
	return rc.output
}

func (rc *stageResource) Error() chan<- error {
	return rc.err
}

type fifo struct {
	proc Processor
}

// FIFO implements the Job interface
func FIFO(proc Processor) Stage {
	return fifo{proc: proc}
}

// process data coming from the input and send out to the output channel after processing
func (f fifo) process(ctx context.Context, input Data, s StageResource) {
	output, err := f.proc.Process(ctx, input)
	if err != nil {
		handleError(err, s.Error())
		return
	}
	if output == nil {
		return
	}

	select {
	case <-ctx.Done():
		log.Printf("Context done %v\n", ctx.Err())
	case s.Output() <- output:
	}
}

// Run the stage infinitely until the context will be eventually canceled
func (f fifo) Run(ctx context.Context, s StageResource) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Ending stage run %v\n", ctx.Err())
			return
		case input, ok := <-s.Input():
			if !ok {
				return
			}
			f.process(ctx, input, s)
		}
	}
}

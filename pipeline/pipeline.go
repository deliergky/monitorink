package pipeline

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Pipeline processes data from a producer through multiple stages and the final processed data
// is consumed by it's consumer
type Pipeline struct {
	stages []Stage
}

// New creates a new pipeline with the provided stages
func New(stages ...Stage) *Pipeline {
	return &Pipeline{
		stages: stages,
	}
}

// Process runs each stage with initial data produces by the Producer and finally consumed by the
// Consumer
func (p *Pipeline) Process(ctx context.Context, prd Producer, cns Consumer) error {
	var wg sync.WaitGroup
	pCtx, cancel := context.WithCancel(ctx)
	// init channels
	rCtxCh := make([]chan Data, len(p.stages)+1)
	errCh := make(chan error, len(p.stages)+1)
	for i := 0; i < len(rCtxCh); i++ {
		rCtxCh[i] = make(chan Data)
	}

	for i := 0; i < len(p.stages); i++ {
		wg.Add(1)
		go func(index int) {
			p.stages[index].Run(pCtx, &stageResource{
				input:  rCtxCh[index],
				output: rCtxCh[index+1],
				err:    errCh,
			})

			close(rCtxCh[index+1])
			wg.Done()
		}(i)
	}

	wg.Add(2)
	go func() {
		produce(pCtx, prd, rCtxCh[0])
		close(rCtxCh[0])
		wg.Done()
	}()

	go func() {
		consume(pCtx, cns, rCtxCh[len(rCtxCh)-1], errCh)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errCh)
		cancel()
	}()
	var wrappedErr error
	for err := range errCh {
		wrappedErr = fmt.Errorf("%w", err)
	}
	cancel()
	return wrappedErr
}

// produce starts an infinite loop to on the Producer which iterates until completion
// the Data is send out on the channel to be processed by stages of the pipeline
func produce(ctx context.Context, prd Producer, out chan<- Data) {
	for prd.Next(ctx) {
		m := prd.Current()
		select {
		case <-ctx.Done():
			log.Printf("Producer got done %v\n", ctx.Err())
			return
		case out <- m:
		}
	}
}

// consume starts an infinite loop to on the Consumer which iterates until completion
// collects the data from the input channel and calls the Consume method implemented
func consume(ctx context.Context, cns Consumer, in <-chan Data, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer got done %v\n", ctx.Err())
			return
		case m, ok := <-in:
			if !ok {
				return
			}
			if err := cns.Consume(ctx, m); err != nil {
				handleError(err, errCh)
			} else {
				m.MarkAsProcessed()
			}
		}
	}
}

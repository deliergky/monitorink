package pipeline

import "log"

func handleError(err error, errCh chan<- error) {
	log.Printf("Handling error %v", err)
	select {
	case errCh <- err:
	default:
	}
}

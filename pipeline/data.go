package pipeline

// Data encapsulates the structure that is being processed
// by the pipeline
type Data interface {
	// MarkAsProcessed identifies that the data processing has been completed
	// at one stage of the pipeline
	MarkAsProcessed()
}

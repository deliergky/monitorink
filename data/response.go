package data

import (
	"context"
	"time"
)

// Iterator is implemented by types that have a collection to iterate on
type Iterator interface {
	Next(context.Context) bool
}

// ResponseData records information of a heartbat/healtcheck reuqest
type ResponseData struct {
	URL          string    `json:"url"`
	ResponseTime int64     `json:"response_time"`
	StatusCode   int       `json:"status_code"`
	Matched      bool      `json:"matched"`
	CreatedAt    time.Time `json:"created_at"`
	processed    bool
}

func (r ResponseData) MarkAsProcessed() {
	r.processed = true
}

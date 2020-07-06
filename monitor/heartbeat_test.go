package monitor

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deliergky/monitorink/data"
	"github.com/deliergky/monitorink/queue"
	"github.com/deliergky/monitorink/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HeartbeatTestSuite struct {
	suite.Suite
	queue  *queue.InMemoryQueue
	store  *store.InMemoryStore
	ctx    context.Context
	cancel context.CancelFunc
	ts     *httptest.Server
}

func (s *HeartbeatTestSuite) SetupTest() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	s.queue = &queue.InMemoryQueue{}
	s.store = &store.InMemoryStore{}

	s.ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "42")
	}))

	for i := 0; i < 10; i++ {
		s.queue.Push(data.ResponseData{})
	}
}

func (s *HeartbeatTestSuite) TestHeartbeat() {
	defer s.ts.Close()
	hbConfig := HeartbeatConfig{
		URL:      s.ts.URL,
		Interval: 1 * time.Second,
	}
	producer := NewResponseProducer(s.queue)
	consumer := NewResponseConsumer(s.store)
	hb := NewHeartbeat(hbConfig, producer, consumer)
	err := hb.Beat(s.ctx)
	assert.NoError(s.T(), err)

	responses, err := s.store.Retrieve(s.ctx)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), responses, 10)
	assert.Equal(s.T(), http.StatusOK, responses[0].StatusCode)

	_, err = s.queue.Pop()
	assert.Error(s.T(), err)
}

func TestStageTestSuite(t *testing.T) {
	suite.Run(t, new(HeartbeatTestSuite))
}

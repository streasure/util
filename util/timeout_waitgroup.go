package util

import (
	"context"
	"errors"
	"sync"
	"time"
)

var errWaitGroupAddFailed = errors.New("timeout wait group: add failed")

type ITimeoutWaitGroup interface {
	Add(delta int) error
	Done()
	WaitTimeout(ctx context.Context, dur time.Duration) bool
}

func NewTimeoutWaitGroup(n int) ITimeoutWaitGroup {
	if n <= 0 {
		return nil
	}
	return &semTimeoutWaitGroup{
		max: n,
	}
}

type semTimeoutWaitGroup struct {
	max   int
	count int
	mu    sync.Mutex
	cond  *sync.Cond
}

func (s *semTimeoutWaitGroup) Add(delta int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count+delta > s.max {
		return errWaitGroupAddFailed
	}
	s.count += delta
	return nil
}

func (s *semTimeoutWaitGroup) Done() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count > 0 {
		s.count--
	}
}

func (s *semTimeoutWaitGroup) WaitTimeout(ctx context.Context, dur time.Duration) bool {
	ctx, cancel := context.WithTimeout(ctx, dur)
	defer cancel()

	for {
		s.mu.Lock()
		if s.count <= 0 {
			s.mu.Unlock()
			return false
		}
		s.mu.Unlock()

		select {
		case <-ctx.Done():
			return true
		case <-time.After(10 * time.Millisecond):
		}
	}
}

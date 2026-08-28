package util

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

type PanicCatcher struct {
	recovered atomic.Pointer[RecoveredPanic]
}

func (p *PanicCatcher) Try(f func()) {
	defer p.tryRecover()
	f()
}

func (p *PanicCatcher) tryRecover() {
	if val := recover(); val != nil {
		rp := NewRecoveredPanic(1, val)
		p.recovered.CompareAndSwap(nil, &rp)
	}
}

func (p *PanicCatcher) Repanic() {
	if val := p.Recovered(); val != nil {
		panic(val)
	}
}

func (p *PanicCatcher) Recovered() *RecoveredPanic {
	return p.recovered.Load()
}

func NewRecoveredPanic(skip int, value any) RecoveredPanic {
	var callers [64]uintptr
	n := runtime.Callers(skip+1, callers[:])
	return RecoveredPanic{
		Value:   value,
		Callers: callers[:n],
		Stack:   debug.Stack(),
	}
}

type RecoveredPanic struct {
	Value   any
	Callers []uintptr
	Stack   []byte
}

func (c *RecoveredPanic) Error() string {
	return fmt.Sprintf("panic: %+v\nstacktrace:\n%s\n", c.Value, c.Stack)
}

func (c *RecoveredPanic) Unwrap() error {
	if err, ok := c.Value.(error); ok {
		return err
	}
	return nil
}
